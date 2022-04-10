package main

import (
	"net/http"

	"github.com/danielblagy/blog-webapp-server/controller"
	"github.com/danielblagy/blog-webapp-server/service"
	"github.com/gin-gonic/gin"
)

var (
	usersService    = service.CreateUsersService()
	usersController = controller.CreateUsersController(usersService)
)

func main() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})

	router.GET("/users", usersController.GetAll)

	router.Run(":4000")
}
