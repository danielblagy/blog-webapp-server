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
	Save(c *gin.Context)
	Unsave(c *gin.Context)
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
	claims, ok := auth.SilentlyCheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	userId := "-1"
	if ok {
		userId = claims.Id
	}

	// if user is unauthorized, userId will be '-1' (used in the service to hide private articles)

	article, err := controller.service.GetById(c.Param("id"), userId)

	if err != nil {
		if err.Error() == "article is private" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

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

	// TODO : test title validation

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
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run update logic

	userId := claims.Id
	articleId := c.Param("id")

	article, err := controller.service.GetById(articleId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	// ensure the user owns the article
	userIdInt, _ := strconv.Atoi(userId)
	if userIdInt != article.AuthorId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "access denied",
		})
		return
	}

	var updatedData entity.EditableArticleData
	if err := c.BindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	updatedArticle, err := controller.service.Update(articleId, updatedData)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedArticle)
}

func (controller *ArticlesControllerProvider) Delete(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run update logic

	userId := claims.Id
	articleId := c.Param("id")

	article, err := controller.service.GetById(articleId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	// ensure the user owns the article
	userIdInt, _ := strconv.Atoi(userId)
	if userIdInt != article.AuthorId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "access denied",
		})
		return
	}

	deletedArticle, err := controller.service.Delete(articleId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, deletedArticle)
}

func (controller *ArticlesControllerProvider) Save(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run 'me' logic

	userId := claims.Id

	// check if the article to save exists
	articleToSave := c.Param("id")
	article, err := controller.service.GetById(articleToSave, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "article to save was not found",
		})
		return
	}

	if err := controller.service.Save(userId, articleToSave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (controller *ArticlesControllerProvider) Unsave(c *gin.Context) {
	claims, ok := auth.CheckForAuthorization(c, "accessToken", "ACCESS_SECRET")
	if !ok {
		return
	}

	// if a token is provided and valid, run 'me' logic

	userId := claims.Id

	// check if the article to unsave exists
	articleToUnsave := c.Param("id")
	article, err := controller.service.GetById(articleToUnsave, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "article to unsave was not found",
		})
		return
	}

	if err := controller.service.Unsave(userId, articleToUnsave); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, article)
}
