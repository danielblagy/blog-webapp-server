package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateJWTToken(userId string, secretKeyEnvVariable string, expirationTime time.Time) (string, error) {
	secretKey := os.Getenv(secretKeyEnvVariable)

	claims := jwt.StandardClaims{}
	//claims["authorized"] = true
	claims.Id = userId
	claims.ExpiresAt = expirationTime.Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(secretKey))
	return token, err
}

func CreateTokenPair(c *gin.Context, userId string) {
	accessTokenDuration := time.Minute * 15
	refreshTokenDuration := time.Hour * 24 * 21

	accessTokenExpirationTime := time.Now().Add(accessTokenDuration)
	refreshTokenExpirationTime := time.Now().Add(refreshTokenDuration)

	accessToken, err := GenerateJWTToken(userId, "ACCESS_SECRET", accessTokenExpirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	refreshToken, err := GenerateJWTToken(userId, "REFRESH_SECRET", refreshTokenExpirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("accessToken", accessToken, int(accessTokenDuration.Seconds()), "/", "", false, false)
	c.SetCookie("refreshToken", refreshToken, int(refreshTokenDuration.Seconds()), "/", "", false, false)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func CheckForAuthorization(c *gin.Context, cookieName string, secretKeyEnvVariable string) (jwt.StandardClaims, bool) {
	tokenString, err := c.Cookie(cookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return jwt.StandardClaims{}, false
	}

	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(secretKeyEnvVariable)), nil
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
