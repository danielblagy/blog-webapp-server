package routes

import (
	"github.com/danielblagy/blog-webapp-server/controller"
	"github.com/gin-gonic/gin"
)

func CreateUsersRoutes(apiGroup *gin.RouterGroup, usersController controller.UsersController) {
	users := apiGroup.Group("/users")

	users.GET("/", usersController.GetAll)
	users.GET("/:id", usersController.GetById)

	users.POST("/signup", usersController.Create)
	users.POST("/signin", usersController.SignIn)
	users.POST("/refresh", usersController.Refresh)

	users.PUT("/", usersController.Update)
	// TODO: create delete /:id endpoint for administrators
	users.DELETE("/", usersController.Delete)
}
