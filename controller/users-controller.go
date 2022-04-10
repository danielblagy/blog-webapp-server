package controller

import (
	"net/http"

	"github.com/danielblagy/blog-webapp-server/service"
	"github.com/gin-gonic/gin"
)

type UsersController interface {
	GetAll(c *gin.Context)
	GetById(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
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

}

func (controller *UsersControllerProvider) Update(c *gin.Context) {

}

func (controller *UsersControllerProvider) Delete(c *gin.Context) {

}
