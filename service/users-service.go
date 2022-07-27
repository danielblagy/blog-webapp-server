package service

import (
	"errors"
	"strconv"

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
	Follow(userId string, userToFollow string) error
	Unfollow(userId string, userToUnfollow string) error
}

type UsersServiceProvider struct {
	database        *gorm.DB
	articlesService ArticlesService
}

func CreateUsersService(database *gorm.DB, articlesService ArticlesService) UsersService {
	return &UsersServiceProvider{
		database:        database,
		articlesService: articlesService,
	}
}

func (service *UsersServiceProvider) loadAssociatedData(user *entity.User, accessLevelCondition string) error {
	if result := service.database.Where(accessLevelCondition, user.Id).Find(&user.Articles); result.Error != nil {
		return errors.New("failed to load associated data")
	}

	// load articles associated data
	for i := range user.Articles {
		if err := service.articlesService.LoadAssociatedData(&user.Articles[i]); err != nil {
			return err
		}
	}

	service.loadAssociatedFollowersData(user)

	return nil
}

// Won't load user.articles and articles associated data
func (service *UsersServiceProvider) loadAssociatedFollowersData(user *entity.User) error {
	var count int64
	// loading user.followers
	if result := service.database.Model(&entity.Follower{}).Where("follows_id = ?", user.Id).Count(&count); result.Error != nil {
		return errors.New("failed to load associated data")
	}
	user.Followers = int(count)
	// loading user.following
	if result := service.database.Model(&entity.Follower{}).Where("follower_id = ?", user.Id).Count(&count); result.Error != nil {
		return errors.New("failed to load associated data")
	}
	user.Following = int(count)

	return nil
}

func (service *UsersServiceProvider) GetAll() ([]entity.User, error) {
	var users []entity.User
	result := service.database.Find(&users)

	// load users associated data
	for i := range users {
		// TODO : handle "failed to load associated data" error in controller
		if err := service.loadAssociatedFollowersData(&users[i]); err != nil {
			return users, result.Error
		}
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

	// TODO : handle "failed to load associated data" error in controller
	if err := service.loadAssociatedData(&user, condition); err != nil {
		return user, result.Error
	}

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

func (service *UsersServiceProvider) Follow(userId string, userToFollow string) error {
	iUserId, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}

	iUserToFollow, err := strconv.Atoi(userToFollow)
	if err != nil {
		return err
	}

	result := service.database.Create(&entity.Follower{FollowerId: iUserId, FollowsId: iUserToFollow})
	return result.Error
}

func (service *UsersServiceProvider) Unfollow(userId string, userToUnfollow string) error {
	result := service.database.Where("follower_id = ? and follows_id = ?", userId, userToUnfollow).Delete(&entity.Follower{})
	return result.Error
}
