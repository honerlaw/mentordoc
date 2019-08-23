package server

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserService struct {
	userRepository        *UserRepository
	organizationService   *OrganizationService
	transactionManager    *util.TransactionManager
}

func NewUserService(
	userRepository *UserRepository,
	organizationService *OrganizationService,
	transactionManager *util.TransactionManager) *UserService {

	service := &UserService{
		userRepository:        userRepository,
		organizationService:   organizationService,
		transactionManager:    transactionManager,
	};
	return service
}

func (service *UserService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewUserService(
		service.userRepository.InjectTransaction(tx).(*UserRepository),
		service.organizationService.InjectTransaction(tx).(*OrganizationService),
		service.transactionManager.InjectTransaction(tx).(*util.TransactionManager))
}

func (service *UserService) Create(email string, password string) (*model.User, error) {
	user := service.userRepository.FindByEmail(email)
	if user != nil {
		return nil, model.NewBadRequestError("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return nil, model.NewInternalServerError("failed to create user")
	}

	resp, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*UserService)

		user = &model.User{
			Email:    email,
			Password: string(hash),
		}
		user.Id = uuid.NewV4().String()

		user, err = injectedService.userRepository.Insert(user)

		if err != nil {
			return nil, model.NewInternalServerError("failed to create user")
		}

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return resp.(*model.User), nil
}

func (service *UserService) Authenticate(email string, password string) (*model.User, error) {
	user := service.userRepository.FindByEmail(email)
	if user == nil {
		return nil, model.NewBadRequestError("invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, model.NewBadRequestError("invalid email or password")
	}

	return user, nil
}

func (service *UserService) FindByEmail(email string) *model.User {
	return service.userRepository.FindByEmail(email)
}
