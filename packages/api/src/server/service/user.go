package service

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
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

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return nil, errors.New("failed to create user")
	}

	user = &model.User{
		Email: email,
		Password: string(hash),
	}
	user.Id = uuid.Must(uuid.NewV4()).String()

	user, err = service.userDao.Insert(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (service *User) Authenticate(email string, password string) (*model.User, error) {
	user := service.userDao.FindByEmail(email)
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (service *User) FindByEmail(email string) *model.User {
	return service.userDao.FindByEmail(email)
}
