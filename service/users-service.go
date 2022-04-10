package service

import (
	"errors"

	"github.com/danielblagy/blog-webapp-server/entity"
)

type UsersService interface {
	GetAll() ([]entity.User, error)
	GetById(id string) (entity.User, error)
	Create(id string, user entity.User) (entity.User, error)
	Update(id string, user entity.User) (entity.User, error)
	Delete(id string) (entity.User, error)
}

type UsersServiceProvider struct {
	users []entity.User
}

func CreateUsersService() UsersService {
	return &UsersServiceProvider{
		users: []entity.User{
			{Id: "1", Login: "userOne", FullName: "Dave", Password: "12345"},
			{Id: "2", Login: "user22", FullName: "Jane", Password: "8888"},
			{Id: "3", Login: "userTres", FullName: "Maggie", Password: "37212asd"},
		},
	}
}

func (service *UsersServiceProvider) findUser(id string) (entity.User, error) {
	for _, user := range service.users {
		if user.Id == id {
			return user, nil
		}
	}

	return entity.User{}, errors.New("user not found")
}

func (service *UsersServiceProvider) findUserIndex(id string) (int, error) {
	for i, user := range service.users {
		if user.Id == id {
			return i, nil
		}
	}

	return -1, errors.New("user not found")
}

func (service *UsersServiceProvider) GetAll() ([]entity.User, error) {
	return service.users, nil
}

func (service *UsersServiceProvider) GetById(id string) (entity.User, error) {
	user, err := service.findUser(id)
	return user, err
}

func (service *UsersServiceProvider) Create(id string, user entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (service *UsersServiceProvider) Update(id string, user entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (service *UsersServiceProvider) Delete(id string) (entity.User, error) {
	return entity.User{}, nil
}
