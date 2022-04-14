package routes

import (
	"github.com/danielblagy/blog-webapp-server/controller"
	"github.com/gin-gonic/gin"
)

func CreateArticlesRoutes(apiGroup *gin.RouterGroup, articlesController controller.ArticlesController) {
	users := apiGroup.Group("/articles")

	users.GET("/", articlesController.GetAll)
	users.GET("/:id", articlesController.GetById)

	users.POST("/", articlesController.Create)

	users.PUT("/:id", articlesController.Update)
	users.DELETE("/:id", articlesController.Delete)
}
