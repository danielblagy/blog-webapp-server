package main

import (
	"log"
	"net/http"

	"github.com/danielblagy/blog-webapp-server/controller"
	"github.com/danielblagy/blog-webapp-server/db"
	"github.com/danielblagy/blog-webapp-server/routes"
	"github.com/danielblagy/blog-webapp-server/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	usersService    service.UsersService
	usersController controller.UsersController

	database          *gorm.DB
	dbConnectionError error
)

func main() {
	database, dbConnectionError = db.SetUpConnection()
	if dbConnectionError != nil {
		log.Fatalf("Failed to set up DB connection: %s", dbConnectionError.Error())
		return
	}

	// init services and controllers

	usersService = service.CreateUsersService(database)
	usersController = controller.CreateUsersController(usersService)

	// set up gin router

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world!",
		})
	})

	api := router.Group("/")
	routes.CreateUsersRoutes(api, usersController)

	router.Run(":4000")
}
