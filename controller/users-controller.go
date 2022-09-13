package controller

import (
	"net/http"
	"strconv"

	"github.com/danielblagy/blog-webapp-server/auth"
	"github.com/danielblagy/blog-webapp-server/entity"
	"github.com/danielblagy/blog-webapp-server/service"
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
	Refresh(c *gin.Context)
	Me(c *gin.Context)
	Follow(c *gin.Context)
	Unfollow(c *gin.Context)
	GetFollowers(c *gin.Context)
	GetFollowing(c *gin.Context)
	IsFollowed(c *gin.Context)
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
	UserToGetId := c.Param("id")
	hasAccessToPrivateArticles := false
	claims, authorizedUser := auth.SilentlyCheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if authorizedUser {
		AuthorizedUserId := claims.Id
		hasAccessToPrivateArticles = AuthorizedUserId == UserToGetId
	}

	user, err := controller.service.GetById(UserToGetId, hasAccessToPrivateArticles)

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

	// TODO: test user data validation

	if newUser.Login == "" || newUser.Password == "" || newUser.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user data",
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
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
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

// TODO: check if client is an authorized administrator
func (controller *UsersControllerProvider) Delete(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run delete logic

	userId := claims.Id

	user, err := controller.service.Delete(userId)

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

	auth.CreateTokenPair(c, strconv.Itoa(user.Id))
}

func (controller *UsersControllerProvider) Refresh(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "refreshToken", "REFRESH_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run update logic

	userId := claims.Id

	auth.CreateTokenPair(c, userId)
}

func (controller *UsersControllerProvider) Me(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run 'me' logic

	userId := claims.Id

	user, err := controller.service.GetById(userId, true)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (controller *UsersControllerProvider) Follow(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run 'me' logic

	userId := claims.Id

	user, err := controller.service.GetById(userId, true)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	// check if the user to follow exists
	userToFollow := c.Param("id")
	_, err = controller.service.GetById(userToFollow, false)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user to follow was not found",
		})
		return
	}

	if err := controller.service.Follow(userId, userToFollow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (controller *UsersControllerProvider) Unfollow(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run 'me' logic

	userId := claims.Id

	user, err := controller.service.GetById(userId, true)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	// check if the user to unfollow exists
	userToUnfollow := c.Param("id")
	_, err = controller.service.GetById(userToUnfollow, false)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user to unfollow was not found",
		})
		return
	}

	if err := controller.service.Unfollow(userId, userToUnfollow); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (controller *UsersControllerProvider) GetFollowers(c *gin.Context) {
	user := c.Param("id")

	followers, err := controller.service.GetFollowers(user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, followers)
}

func (controller *UsersControllerProvider) GetFollowing(c *gin.Context) {
	user := c.Param("id")

	following, err := controller.service.GetFollowing(user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, following)
}

func (controller *UsersControllerProvider) IsFollowed(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run logic

	userId := claims.Id

	userToCheck := c.Param("id")
	isFollowed, err := controller.service.IsFollowed(userId, userToCheck)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, isFollowed)
}
