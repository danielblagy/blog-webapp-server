package controller

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/danielblagy/blog-webapp-server/entity"
	"github.com/danielblagy/blog-webapp-server/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UsersController interface {
	GetAll(c *gin.Context)
	GetById(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	SignIn(c *gin.Context)
}

type UsersControllerProvider struct {
	service service.UsersService
}

func CreateUsersController(service service.UsersService) UsersController {
	return &UsersControllerProvider{
		service: service,
	}
}

func (controller *UsersControllerProvider) GetAll(c *gin.Context) {
	users, err := controller.service.GetAll()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (controller *UsersControllerProvider) GetById(c *gin.Context) {
	user, err := controller.service.GetById(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (controller *UsersControllerProvider) Create(c *gin.Context) {
	var newUser entity.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	_, err := controller.service.GetByLogin(newUser.Login)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "this login is taken",
		})
		return
	}

	createdUser, err := controller.service.Create(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (controller *UsersControllerProvider) Update(c *gin.Context) {
	claims, ok := controller.checkForAuthorization(c)
	if !ok {
		return
	}

	// if a token is provided and valid, run update logic

	userId := claims.Id

	var updatedData entity.EditableUserData
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	updatedUser, err := controller.service.Update(userId, updatedData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (controller *UsersControllerProvider) Delete(c *gin.Context) {
	user, err := controller.service.Delete(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (controller *UsersControllerProvider) SignIn(c *gin.Context) {
	var claimedUser entity.User
	if err := c.BindJSON(&claimedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := controller.service.GetByLogin(claimedUser.Login)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user with this login doesn't exist",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(claimedUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	expirationTime := time.Now().Add(time.Minute * 15)

	token, err := controller.generateJWTToken(user, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("token", token, int((time.Minute * 15).Seconds()), "/", "", false, false)
	/*http.SetCookie(c.Writer, &http.Cookie{
		Name:    "jws_token",
		Value:   token,
		Expires: expirationTime,
	})*/
	c.JSON(http.StatusOK, token)
}

func (controller *UsersControllerProvider) generateJWTToken(user entity.User, expirationTime time.Time) (string, error) {
	secretKey := os.Getenv("ACCESS_SECRET")

	claims := jwt.StandardClaims{}
	//claims["authorized"] = true
	claims.Id = strconv.Itoa(user.Id)
	claims.ExpiresAt = expirationTime.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(secretKey))
	return token, err
}

func (controller *UsersControllerProvider) checkForAuthorization(c *gin.Context) (jwt.StandardClaims, bool) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return jwt.StandardClaims{}, false
	}

	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return jwt.StandardClaims{}, false
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return jwt.StandardClaims{}, false
		}
	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "access denied",
		})
		return jwt.StandardClaims{}, false
	}

	return claims, true
}
