package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
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

func (service *AclService) LinkUserToRole(user *shared.User, roleName string, resourceId string) error {
	return service.userRoleService.LinkUserToRole(user, roleName, resourceId)
}

func (service *AclService) UserCanAccessResourceByModel(user *shared.User, model interface{}, action string) bool {
	data, err := service.GetResourceDataForModel(model)
	if err != nil {
		log.Print(err)
		return false
	}
	return service.UserCanAccessResource(user, data.ResourcePath, data.ResourceIds, action)
}

func (service *AclService) UserCanAccessResource(user *shared.User, path []string, ids []string, action string) bool {
	canAccess, err := service.userRoleService.UserCanAccessResource(user, path, ids, action)
	if err != nil {
		log.Print(err)
		return false
	}
	return canAccess
}

func (service *AclService) UserActionableResourcesByPath(user *shared.User, path []string, action string) ([]ResourceResponse, error) {
	return service.userRoleService.UserActionableResourcesByPath(user, path, action)
}

func (service *AclService) UserActionsForResources(user *shared.User, paths [][]string, ids [][]string) ([]ResourceResponse, error) {
	return service.userRoleService.UserActionsForResources(user, paths, ids)
}

func (service *AclService) Wrap(user *shared.User, modelSlice interface{}) ([]AclWrappedModel, error) {
	return service.aclWrapperService.Wrap(user, modelSlice)
}

func (service *AclService) GetResourceDataForModel(model interface{}) (*ResourceData, error) {
	return service.aclWrapperService.GetResourceDataForModel(model)
}
