package service

import "github.com/danielblagy/blog-webapp-server/entity"

type UsersService interface {
	GetAll() ([]entity.User, error)
	GetById(id int) (entity.User, error)
	Create(id int, user entity.User) (entity.User, error)
	Update(id int, user entity.User) (entity.User, error)
	Delete(id int) (entity.User, error)
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

func (service UsersServiceProvider) GetAll() ([]entity.User, error) {
	return service.users, nil
}

func (service UsersServiceProvider) GetById(id int) (entity.User, error) {
	return entity.User{}, nil
}

func (service UsersServiceProvider) Create(id int, user entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (service UsersServiceProvider) Update(id int, user entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (service UsersServiceProvider) Delete(id int) (entity.User, error) {
	return entity.User{}, nil
}
