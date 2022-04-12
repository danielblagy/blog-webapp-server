package service

import (
	"github.com/danielblagy/blog-webapp-server/entity"
	"gorm.io/gorm"
)

type UsersService interface {
	GetAll() ([]entity.User, error)
	GetById(id string) (entity.User, error)
	Create(user entity.User) (entity.User, error)
	Update(id string, user entity.User) (entity.User, error)
	Delete(id string) (entity.User, error)
}

type UsersServiceProvider struct {
	database *gorm.DB
	users    []entity.User
}

func CreateUsersService(database *gorm.DB) UsersService {
	return &UsersServiceProvider{
		database: database,
		users:    []entity.User{},
		/*users: []entity.User{
			{Id: 1, Login: "userOne", FullName: "Dave", Password: "12345", Articles: []entity.Article{}},
			{Id: 2, Login: "user22", FullName: "Jane", Password: "8888", Articles: []entity.Article{}},
			{Id: 3, Login: "userTres", FullName: "Maggie", Password: "37212asd", Articles: []entity.Article{}},
		},*/
	}
}

func (service *UsersServiceProvider) findUser(idStr string) (entity.User, error) {
	/*id, err := strconv.Atoi(idStr)
	if err != nil {
		return entity.User{}, errors.New("server error")
	}

	for _, user := range service.users {
		if user.Id == id {
			return user, nil
		}
	}

	return entity.User{}, errors.New("user not found")*/
	return entity.User{}, nil
}

func (service *UsersServiceProvider) findUserIndex(idStr string) (int, error) {
	/*id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, errors.New("server error")
	}

	for i, user := range service.users {
		if user.Id == id {
			return i, nil
		}
	}

	return -1, errors.New("user not found")*/
	return 0, nil
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

func (service *UsersServiceProvider) Create(user entity.User) (entity.User, error) {
	result := service.database.Create(&user)
	return user, result.Error
}

func (service *UsersServiceProvider) Update(id string, user entity.User) (entity.User, error) {
	/*i, err := service.findUserIndex(id)

	if err == nil {
		service.users[i] = user
	}

	return user, err*/

	result := service.database.Save(&user)
	return user, result.Error
}

func (service *UsersServiceProvider) Delete(id string) (entity.User, error) {
	user, _ := service.GetById(id)
	result := service.database.Delete(&entity.User{}, id)
	return user, result.Error
}
