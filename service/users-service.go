package service

import (
	"github.com/danielblagy/blog-webapp-server/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersService interface {
	GetAll() ([]entity.User, error)
	GetById(id string, authorized bool) (entity.User, error)
	GetByLogin(login string) (entity.User, error)
	Create(user entity.User) (entity.User, error)
	Update(id string, updatedData entity.EditableUserData) (entity.User, error)
	Delete(id string) (entity.User, error)
}

type UsersServiceProvider struct {
	database *gorm.DB
}

func CreateUsersService(database *gorm.DB) UsersService {
	return &UsersServiceProvider{
		database: database,
	}
}

func (service *UsersServiceProvider) GetAll() ([]entity.User, error) {
	var users []entity.User
	result := service.database.Find(&users)

	// TODO: deal with getting associated data
	// associated data
	var count int64
	for i := range users {
		// set followers count field
		service.database.Model(&entity.Follower{}).Where("follows_id = ?", users[i].Id).Count(&count)
		users[i].Followers = int(count)
		// set following count field
		service.database.Model(&entity.Follower{}).Where("follower_id = ?", users[i].Id).Count(&count)
		users[i].Following = int(count)
	}

	return users, result.Error
}

func (service *UsersServiceProvider) GetById(id string, authorized bool) (entity.User, error) {
	var user entity.User
	result := service.database.First(&user, id)

	condition := "author_id = ?"
	if !authorized {
		condition += " and published = true"
	}

	// associated data
	// TODO: deal with this mess
	service.database.Where(condition, user.Id).Find(&user.Articles)
	// set article.author field for every article
	for i := range user.Articles {
		service.database.Where("id = ?", user.Id).Find(&user.Articles[i].Author)
	}
	// set followers count field
	var count int64
	service.database.Model(&entity.Follower{}).Where("follows_id = ?", user.Id).Count(&count)
	user.Followers = int(count)
	// set following count field
	service.database.Model(&entity.Follower{}).Where("follower_id = ?", user.Id).Count(&count)
	user.Following = int(count)

	return user, result.Error
}

func (service *UsersServiceProvider) GetByLogin(login string) (entity.User, error) {
	var user entity.User
	result := service.database.Where("login = ?", login).First(&user)
	return user, result.Error
}

func (service *UsersServiceProvider) Create(user entity.User) (entity.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
	if err != nil {
		return user, err
	}
	user.Password = string(hash)

	result := service.database.Create(&user)
	return user, result.Error
}

func (service *UsersServiceProvider) Update(id string, updatedData entity.EditableUserData) (entity.User, error) {
	var user entity.User
	service.database.Find(&user, id)

	if updatedData.FullName != "" {
		user.FullName = updatedData.FullName
	}

	if updatedData.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(updatedData.Password), 0)
		if err != nil {
			return user, err
		}
		user.Password = string(hash)
	}

	result := service.database.Save(&user)
	return user, result.Error
}

func (service *UsersServiceProvider) Delete(id string) (entity.User, error) {
	user, _ := service.GetById(id, true) // getting the user before deleting to return
	result := service.database.Delete(&entity.User{}, id)
	return user, result.Error
}
