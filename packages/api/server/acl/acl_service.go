package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
)

type AclService struct {
	rolePermissionService *RolePermissionService
	userRoleService       *UserRoleService
	transactionManager    *server.TransactionManager
}

func NewAclService(rolePermissionService *RolePermissionService,
	userRoleService *UserRoleService,
	transactionManager *server.TransactionManager) *AclService {
	return &AclService{
		rolePermissionService: rolePermissionService,
		userRoleService:       userRoleService,
		transactionManager:    transactionManager,
	}
}

func (service *AclService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewAclService(service.rolePermissionService.InjectTransaction(tx).(*RolePermissionService),
		service.userRoleService.InjectTransaction(tx).(*UserRoleService),
		service.transactionManager.InjectTransaction(tx).(*server.TransactionManager))
}

func (service *AclService) Init() error {
	return service.rolePermissionService.InitRoles()
}

func (service *AclService) LinkUserToRole(user model.User, roleName string) error {
	return service.userRoleService.LinkUserToRole(user, roleName)
}

func (service *AclService) GetRolesForUser(user model.User) error {

}