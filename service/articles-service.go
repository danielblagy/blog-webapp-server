package service

import (
	"errors"
	"strconv"

	"github.com/danielblagy/blog-webapp-server/entity"
	"gorm.io/gorm"
)

type ArticlesService interface {
	LoadAssociatedData(*entity.Article) error
	GetAll() ([]entity.Article, error)
	GetById(id string, userId string) (entity.Article, error)
	GetByTitle(authorId string, title string) (entity.Article, error)
	Create(article entity.Article) (entity.Article, error)
	Update(id string, updatedData entity.EditableArticleData) (entity.Article, error)
	Delete(id string) (entity.Article, error)
	Save(userId string, articleToSave string) error
	Unsave(userId string, articleToUnsave string) error
	GetSaves(userId string) ([]entity.Article, error)
	IsSaved(userId string, articleId string) (bool, error)
	ForYou(userId string) ([]entity.Article, error)
}

type ArticlesServiceProvider struct {
	database *gorm.DB
}

func CreateArticlesService(database *gorm.DB) ArticlesService {
	return &ArticlesServiceProvider{
		database: database,
	}
}

func (service *ArticlesServiceProvider) LoadAssociatedData(article *entity.Article) error {
	// NOTE: user's associeated data will not be loaded
	// loading article.author
	result := service.database.Where("id = ?", article.AuthorId).First(&article.Author)
	if result.Error != nil {
		return errors.New("failed to load associated data")
	}

	// loading article.saves
	var count int64
	result = service.database.Model(entity.Save{}).Where("article_id = ?", article.Id).Count(&count)
	if result.Error != nil {
		return errors.New("failed to load associated data")
	}
	article.Saves = int(count)

	return nil
}

func (service *ArticlesServiceProvider) GetAll() ([]entity.Article, error) {
	var articles []entity.Article
	result := service.database.Where("published = true").Find(&articles)

	// associated data
	for i := range articles {
		if err := service.LoadAssociatedData(&articles[i]); err != nil {
			return articles, err
		}
	}

	return articles, result.Error
}

func (service *ArticlesServiceProvider) GetById(id string, userId string) (entity.Article, error) {
	var article entity.Article

	if userId == "-1" {
		result := service.database.Where("published = true").First(&article, id)
		service.database.Where("id = ?", article.AuthorId).Find(&article.Author)
		return article, result.Error
	}

	result := service.database.First(&article, id)

	if article.Published == false && strconv.Itoa(article.AuthorId) != userId {
		return entity.Article{}, errors.New("article is private")
	}

	if err := service.LoadAssociatedData(&article); err != nil {
		return article, err
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

	if err := service.LoadAssociatedData(&article); err != nil {
		return article, err
	}

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

	if err := service.LoadAssociatedData(&article); err != nil {
		return article, err
	}

	return article, result.Error
}

func (service *ArticlesServiceProvider) Delete(id string) (entity.Article, error) {
	//article, _ := service.GetById(id)
	// getting the article before deleting to return
	var article entity.Article
	service.database.First(&article, id)

	if err := service.LoadAssociatedData(&article); err != nil {
		return article, err
	}

	result := service.database.Delete(&entity.Article{}, id)
	return article, result.Error
}

func (service *ArticlesServiceProvider) Save(userId string, articleToSave string) error {
	iUserId, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}

	iArticleToSave, err := strconv.Atoi(articleToSave)
	if err != nil {
		return err
	}

	result := service.database.Create(&entity.Save{UserId: iUserId, ArticleId: iArticleToSave})
	return result.Error
}

func (service *ArticlesServiceProvider) Unsave(userId string, articleToUnsave string) error {
	result := service.database.Where("user_id = ? and article_id = ?", userId, articleToUnsave).Delete(&entity.Save{})
	return result.Error
}

func (service *ArticlesServiceProvider) GetSaves(userId string) ([]entity.Article, error) {
	var savedArticlesIds []int
	result := service.database.Table("saves").Where("user_id = ?", userId).Select("article_id").Find(&savedArticlesIds)

	var articles []entity.Article
	result = service.database.Where("id in ? and published = true", savedArticlesIds).Find(&articles)

	// associated data
	for i := range articles {
		if err := service.LoadAssociatedData(&articles[i]); err != nil {
			return articles, err
		}
	}

	return articles, result.Error
}

func (service *ArticlesServiceProvider) IsSaved(userId string, articleId string) (bool, error) {
	var savedArticlesIds []int
	result := service.database.Table("saves").Where("user_id = ? and article_id = ?", userId, articleId).Select("article_id").Find(&savedArticlesIds)

	return result.RowsAffected > 0, result.Error
}

func (service *ArticlesServiceProvider) ForYou(userId string) ([]entity.Article, error) {
	var following []int
	result := service.database.Table("followers").Where("follower_id = ?", userId).Select("follows_id").Find(&following)

	var articles []entity.Article
	result = service.database.Where("author_id in ? and published = true", following).Find(&articles)

	// associated data
	for i := range articles {
		if err := service.LoadAssociatedData(&articles[i]); err != nil {
			return articles, err
		}
	}

	return articles, result.Error
}
