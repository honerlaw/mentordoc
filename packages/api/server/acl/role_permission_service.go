package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server"
	"github.com/honerlaw/mentordoc/server/model"
	uuid "github.com/satori/go.uuid"
)

type RolePermissionService struct {
	roleRepository           *RoleRepository
	permissionRepository     *PermissionRepository
	rolePermissionRepository *RolePermissionRepository
	transactionManager       *server.TransactionManager
}

func NewRolePermissionService(roleRepository *RoleRepository,
	permissionRepository *PermissionRepository,
	rolePermissionRepository *RolePermissionRepository,
	transactionManager *server.TransactionManager) *RolePermissionService {

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
		service.transactionManager.InjectTransaction(tx).(*server.TransactionManager))
}

func (service *RolePermissionService) InitRoles() error {
	_, err := service.CreateRoleWithPermissions("organization:owner", map[string][]string {
		"organization": {"view", "modify", "create:folder"},
		"organization:folder": {"view", "modify", "delete", "create:document"},
		"organization:document": {"view", "modify", "delete"},
	});
	if err != nil {
		return err
	}
	_, err = service.CreateRoleWithPermissions("organization:contributor", map[string][]string {
		"organization": {"view", "create:folder"},
		"organization:folder": {"view", "modify", "delete", "create:document"},
		"organization:document": {"view", "modify", "delete"},
	});
	if err != nil {
		return err
	}
	return nil
}

func (service *RolePermissionService) CreateRoleWithPermissions(roleName string, permissionMap map[string][]string) (*model.Role, error) {
	role, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*RolePermissionService)

		role := &model.Role{}
		role.Id = uuid.NewV4().String()

		role, err := injectedService.roleRepository.Insert(role)
		if err != nil {
			return nil, err
		}

		for path, actions := range permissionMap {
			for _, action := range actions {
				permission := &model.Permission{
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

	return role.(*model.Role), err
}