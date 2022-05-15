package service

import (
	"errors"
	"strconv"

	"github.com/danielblagy/blog-webapp-server/entity"
	"gorm.io/gorm"
)

type ArticlesService interface {
	GetAll() ([]entity.Article, error)
	GetById(id string, userId string) (entity.Article, error)
	GetByTitle(authorId string, title string) (entity.Article, error)
	Create(article entity.Article) (entity.Article, error)
	Update(id string, updatedData entity.EditableArticleData) (entity.Article, error)
	Delete(id string) (entity.Article, error)
}

type ArticlesServiceProvider struct {
	database *gorm.DB
}

func CreateArticlesService(database *gorm.DB) ArticlesService {
	return &ArticlesServiceProvider{
		database: database,
	}
}

func (service *ArticlesServiceProvider) GetAll() ([]entity.Article, error) {
	var articles []entity.Article
	result := service.database.Where("published = true").Find(&articles)
	return articles, result.Error
}

func (service *ArticlesServiceProvider) GetById(id string, userId string) (entity.Article, error) {
	var article entity.Article

	if userId == "-1" {
		result := service.database.Where("published = true").First(&article, id)
		return article, result.Error
	}

	result := service.database.First(&article, id)

	if article.Published == false && strconv.Itoa(article.AuthorId) != userId {
		return entity.Article{}, errors.New("article is private")
	}

	return article, result.Error
}

func (service *ArticlesServiceProvider) GetByTitle(authorId string, title string) (entity.Article, error) {
	var article entity.Article
	result := service.database.Where("author_id = ? and title = ?", authorId, title).First(&article)
	return article, result.Error
}

func (service *ArticlesServiceProvider) Create(article entity.Article) (entity.Article, error) {
	result := service.database.Create(&article)
	return article, result.Error
}

func (service *ArticlesServiceProvider) Update(id string, updatedData entity.EditableArticleData) (entity.Article, error) {
	var article entity.Article
	service.database.Find(&article, id)

	if updatedData.Title != "" {
		article.Title = updatedData.Title
	}

	if updatedData.Content != article.Content {
		article.Content = updatedData.Content
	}

	if updatedData.Published != article.Published {
		article.Published = updatedData.Published
	}

	result := service.database.Save(&article)
	return article, result.Error
}

func (service *ArticlesServiceProvider) Delete(id string) (entity.Article, error) {
	//article, _ := service.GetById(id)
	// getting the article before deleting to return
	var article entity.Article
	service.database.First(&article, id)

	result := service.database.Delete(&entity.Article{}, id)
	return article, result.Error
}
