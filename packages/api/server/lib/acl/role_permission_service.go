package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/lib/util"
	uuid "github.com/satori/go.uuid"
)

type RolePermissionService struct {
	roleRepository           *RoleRepository
	permissionRepository     *PermissionRepository
	rolePermissionRepository *RolePermissionRepository
	transactionManager       *util.TransactionManager
}

func NewRolePermissionService(roleRepository *RoleRepository,
	permissionRepository *PermissionRepository,
	rolePermissionRepository *RolePermissionRepository,
	transactionManager *util.TransactionManager) *RolePermissionService {

	return &RolePermissionService{
		roleRepository:           roleRepository,
		permissionRepository:     permissionRepository,
		rolePermissionRepository: rolePermissionRepository,
		transactionManager:       transactionManager,
	}
}

func (service *RolePermissionService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewRolePermissionService(service.roleRepository.InjectTransaction(tx).(*RoleRepository),
		service.permissionRepository.InjectTransaction(tx).(*PermissionRepository),
		service.rolePermissionRepository.InjectTransaction(tx).(*RolePermissionRepository),
		service.transactionManager.InjectTransaction(tx).(*util.TransactionManager))
}

func (service *RolePermissionService) InitRoles() error {
	_, err := service.CreateRoleWithPermissions("organization:owner", map[string][]string {
		"organization": {"view", "modify", "view:folder", "create:folder", "create:document", "view:document"},
		"organization:folder": {"view", "modify", "delete", "view:folder", "create:folder", "view:document", "create:document"},
		"organization:folder:document": {"view", "modify", "delete"},
	});
	if err != nil {
		return err
	}
	_, err = service.CreateRoleWithPermissions("organization:contributor", map[string][]string {
		"organization": {"view", "create:folder", "create:document", "view:document"},
		"organization:folder": {"view", "modify", "delete", "view:folder", "create:folder", "view:document", "create:document"},
		"organization:folder:document": {"view", "modify", "delete"},
	});
	if err != nil {
		return err
	}
	return nil
}

func (service *RolePermissionService) CreateRoleWithPermissions(roleName string, permissionMap map[string][]string) (*Role, error) {
	role, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*RolePermissionService)

		role := &Role{
			Name: roleName,
		}
		role.Id = uuid.NewV4().String()

		role, err := injectedService.roleRepository.Insert(role)
		if err != nil {
			return nil, err
		}

		for path, actions := range permissionMap {
			for _, action := range actions {
				permission := &Permission{
					ResourcePath: path,
					Action:       action,
				}
				permission.Id = uuid.NewV4().String()

				permission, err := injectedService.permissionRepository.Insert(permission)
				if err != nil {
					return nil, err
				}
				if err := injectedService.rolePermissionRepository.Link(role, permission); err != nil {
					return nil, err
				}
			}
		}

		return role, nil
	})

	if err != nil {
		return nil, err
	}

	return role.(*Role), nil
}