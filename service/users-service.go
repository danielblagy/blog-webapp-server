package service

import (
	"github.com/danielblagy/blog-webapp-server/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersService interface {
	GetAll() ([]entity.User, error)
	GetById(id string) (entity.User, error)
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
	return users, result.Error
}

func (service *UsersServiceProvider) GetById(id string) (entity.User, error) {
	var user entity.User
	result := service.database.Find(&user, id)
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
	user, _ := service.GetById(id) // getting the user before deleting to return
	result := service.database.Delete(&entity.User{}, id)
	return user, result.Error
}
