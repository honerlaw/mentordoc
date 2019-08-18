package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server"
)

type AclService struct {
	rolePermissionService *RolePermissionService
	transactionManager    *server.TransactionManager
}

func NewAclService(rolePermissionService *RolePermissionService, transactionManager *server.TransactionManager) *AclService {
	return &AclService{
		rolePermissionService: rolePermissionService,
		transactionManager:    transactionManager,
	}
}

func (service *AclService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewAclService(service.rolePermissionService.InjectTransaction(tx).(*RolePermissionService),
		service.transactionManager.InjectTransaction(tx).(*server.TransactionManager))
}

func (service *AclService) Init() error {
	return service.rolePermissionService.InitRoles()
}
