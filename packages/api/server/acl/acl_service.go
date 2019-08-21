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
	db *sql.DB
	tx *sql.Tx
}

/**
ACL should only be accessible through this object and its exposed functions
 */
func NewAclService(transactionManager *server.TransactionManager, db *sql.DB, tx *sql.Tx) *AclService {
	// setup repositories
	roleRepository := NewRoleRepository(db, tx)
	rolePermissionRepository := NewRolePermissionRepository(db, tx)
	permissionRepository := NewPermissionRepository(db, tx)
	userRoleRepository := NewUserRoleRepository(db, tx)

	// setup the permissions
	rolePermissionService := NewRolePermissionService(roleRepository, permissionRepository, rolePermissionRepository, transactionManager)
	userRoleService := NewUserRoleService(roleRepository, userRoleRepository)

	return &AclService{
		rolePermissionService: rolePermissionService,
		userRoleService:       userRoleService,
		transactionManager:    transactionManager,
	}
}

func (service *AclService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewAclService(service.transactionManager.InjectTransaction(tx).(*server.TransactionManager), service.db, tx)
}

func (service *AclService) Init() error {
	return service.rolePermissionService.InitRoles()
}

func (service *AclService) LinkUserToRole(user *model.User, roleName string, resourceId string) error {
	return service.userRoleService.LinkUserToRole(user, roleName, resourceId)
}

func (service *AclService) UserCanAccessResource(user *model.User, path []string, ids []string, action string) (bool, error) {
	return service.userRoleService.UserCanAccessResource(user, path, ids, action)
}

func (service *AclService) UserActionableResourcesByPath(user *model.User, path []string, action string) ([]*ResourceResponse, error) {
	return service.userRoleService.UserActionableResourcesByPath(user, path, action)
}
