package server

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	};
}

func (service *UserService) Create(email string, password string) (*User, error) {
	user := service.userRepository.FindByEmail(email)
	if user != nil {
		return nil, NewBadRequestError("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return nil, NewInternalServerError("failed to create user")
	}

	user = &User{
		Email: email,
		Password: string(hash),
	}
	user.Id = uuid.NewV4().String()

	user, err = service.userRepository.Insert(user)
	if err != nil {
		return nil, NewInternalServerError("failed to create user")
	}

	return user, nil
}

func (service *UserService) Authenticate(email string, password string) (*User, error) {
	user := service.userRepository.FindByEmail(email)
	if user == nil {
		return nil, NewBadRequestError("invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, NewBadRequestError("invalid email or password")
	}

	return user, nil
}

func (service *UserService) FindByEmail(email string) *User {
	return service.userRepository.FindByEmail(email)
}
