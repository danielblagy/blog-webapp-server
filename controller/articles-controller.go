package controller

import (
	"net/http"
	"strconv"

	"github.com/danielblagy/blog-webapp-server/auth"
	"github.com/danielblagy/blog-webapp-server/entity"
	"github.com/danielblagy/blog-webapp-server/service"
	"github.com/gin-gonic/gin"
)

type ArticlesController interface {
	GetAll(c *gin.Context)
	GetById(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type ArticlesControllerProvider struct {
	service service.ArticlesService
}

func CreateArticlesController(service service.ArticlesService) ArticlesController {
	return &ArticlesControllerProvider{
		service: service,
	}
}

func (controller *ArticlesControllerProvider) GetAll(c *gin.Context) {
	articles, err := controller.service.GetAll()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (controller *ArticlesControllerProvider) GetById(c *gin.Context) {
	article, err := controller.service.GetById(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (controller *ArticlesControllerProvider) Create(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run update logic

	userId := claims.Id

	var newArticle entity.Article
	if err := c.BindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// TODO: multiple options: 	1) leave it at that
	//							2) let client set article's author_id field and check if it matches the one in the accessToken
	//							3) use EditableArticleData
	newArticle.AuthorId, _ = strconv.Atoi(userId)

	_, err := controller.service.GetByTitle(strconv.Itoa(newArticle.AuthorId), newArticle.Title)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"message": "user already has article with this title",
		})
		return
	}

	createdArticle, err := controller.service.Create(newArticle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdArticle)
}

func (controller *ArticlesControllerProvider) Update(c *gin.Context) {

}

func (controller *ArticlesControllerProvider) Delete(c *gin.Context) {

}
