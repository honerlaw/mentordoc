package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
)

type AclService struct {
	rolePermissionService *RolePermissionService
	userRoleService       *UserRoleService
	transactionManager    *util.TransactionManager
	aclWrapperService     *AclWrapperService
	db                    *sql.DB
	tx                    *sql.Tx
}

/**
ACL should only be accessible through this object and its exposed functions
*/
func NewAclService(transactionManager *util.TransactionManager, db *sql.DB, tx *sql.Tx) *AclService {
	// setup repositories
	roleRepository := NewRoleRepository(db, tx)
	rolePermissionRepository := NewRolePermissionRepository(db, tx)
	permissionRepository := NewPermissionRepository(db, tx)
	userRoleRepository := NewUserRoleRepository(db, tx)

	// setup the permissions
	rolePermissionService := NewRolePermissionService(roleRepository, permissionRepository, rolePermissionRepository, transactionManager)
	userRoleService := NewUserRoleService(roleRepository, userRoleRepository)

	aclService := &AclService{
		rolePermissionService: rolePermissionService,
		userRoleService:       userRoleService,
		transactionManager:    transactionManager,
		db:                    db,
		tx:                    tx,
	}

	aclService.aclWrapperService = NewAclWrapperService(aclService)

	return aclService
}

func (service *AclService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewAclService(service.transactionManager.InjectTransaction(tx).(*util.TransactionManager), service.db, tx)
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

func (service *AclService) UserActionableResourcesByPath(user *model.User, path []string, action string) ([]ResourceResponse, error) {
	return service.userRoleService.UserActionableResourcesByPath(user, path, action)
}

func (service *AclService) UserActionsForResources(user *model.User, paths [][]string, ids [][]string) ([]ResourceResponse, error) {
	return service.userRoleService.UserActionsForResources(user, paths, ids)
}

func (service *AclService) Wrap(user *model.User, modelSlice interface{}) ([]model.AclWrappedModel, error) {
	return service.aclWrapperService.Wrap(user, modelSlice)
}

func (service *AclService) GetResourceDataForModel(model interface{}) (*ResourceData, error) {
	return service.aclWrapperService.GetResourceDataForModel(model)
}
