package user

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/organization"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

type UserService struct {
	userRepository      *UserRepository
	organizationService *organization.OrganizationService
	transactionManager  *util.TransactionManager
	aclService          *acl.AclService
}

func NewUserService(
	userRepository *UserRepository,
	organizationService *organization.OrganizationService,
	transactionManager *util.TransactionManager,
	aclService *acl.AclService,
) *UserService {

	service := &UserService{
		userRepository:      userRepository,
		organizationService: organizationService,
		transactionManager:  transactionManager,
		aclService:          aclService,
	};
	return service
}

func (service *UserService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewUserService(
		service.userRepository.InjectTransaction(tx).(*UserRepository),
		service.organizationService.InjectTransaction(tx).(*organization.OrganizationService),
		service.transactionManager.InjectTransaction(tx).(*util.TransactionManager),
		service.aclService.InjectTransaction(tx).(*acl.AclService))
}

func (service *UserService) Create(email string, password string) (*shared.User, error) {
	user := service.userRepository.FindByEmail(email)
	if user != nil {
		return nil, shared.NewBadRequestError("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return nil, shared.NewInternalServerError("failed to create user")
	}

	resp, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*UserService)

		user = &shared.User{
			Email:    email,
			Password: string(hash),
		}
		user.Id = uuid.NewV4().String()

		user, err = injectedService.userRepository.Insert(user)
		if err != nil {
			return nil, shared.NewInternalServerError("failed to create user")
		}

		org, err := injectedService.organizationService.Create(strings.Split(user.Email, "@")[0])
		if err != nil {
			return nil, shared.NewInternalServerError("failed to create user")
		}

		err = injectedService.aclService.LinkUserToRole(user, "organization:owner", org.Id)
		if err != nil {
			return nil, shared.NewInternalServerError("failed to create user")
		}

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return resp.(*shared.User), nil
}

func (service *UserService) Authenticate(email string, password string) (*shared.User, error) {
	user := service.userRepository.FindByEmail(email)
	if user == nil {
		return nil, shared.NewBadRequestError("invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, shared.NewBadRequestError("invalid email or password")
	}

	return user, nil
}

func (service *UserService) FindByEmail(email string) *shared.User {
	return service.userRepository.FindByEmail(email)
}

func (service *UserService) FindById(id string) *shared.User {
	return service.userRepository.FindById(id)
}
