package service

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"server/model"
)

type User struct {
	userDao *UserDao
}

func NewUser(userDao *UserDao) *User {
	return &User {
		userDao: userDao,
	};
}

func (service *User) Create(email string, password string) (*model.User, error) {
	user := service.userDao.FindByEmail(email)
	if user == nil {
		return nil, errors.New("failed to create user")
	}

	user = &model.User{
		Email: email,
	}
	user.Id = uuid.Must(uuid.NewV4()).String()

	user, err := service.userDao.Insert(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (service *User) FindByEmail(email string) *model.User {
	return service.userDao.FindByEmail(email)
}
