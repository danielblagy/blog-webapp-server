package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danielblagy/blog-webapp-server/controller"
	"github.com/danielblagy/blog-webapp-server/entity"
	"github.com/danielblagy/blog-webapp-server/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	usersService    service.UsersService
	usersController controller.UsersController

	database          *gorm.DB
	dbConnectionError error
)

func main() {
	// setting up db connection

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	database, dbConnectionError = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if dbConnectionError != nil {
		log.Fatal(dbConnectionError)
	}

	// make migrations to the db (will be done only once, if the entities have never been created before)
	database.AutoMigrate(&entity.User{})
	database.AutoMigrate(&entity.Article{})

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

	// TODO: create a 'users' group

	router.GET("/users", usersController.GetAll)
	router.GET("/users/:id", usersController.GetById)

	router.POST("/users/signup", usersController.Create)
	router.POST("/users/signin", usersController.SignIn)
	router.POST("/users/refresh", usersController.Refresh)

	router.PUT("/users", usersController.Update)
	router.DELETE("/users/:id", usersController.Delete)

	router.Run(":4000")
}
