package acl

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server"
	uuid "github.com/satori/go.uuid"
)

/*
TODO
- check if user has access to do a given action on a resource
- policy templates in general
*/

type AclService struct {
	aclRepository      *AclRepository
	transactionManager *server.TransactionManager
}

func NewAclService(aclRepository *AclRepository, transactionManager *server.TransactionManager) *AclService {
	return &AclService{
		aclRepository:      aclRepository,
		transactionManager: transactionManager,
	}
}

func (service *AclService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewAclService(service.aclRepository.InjectTransaction(tx).(*AclRepository),
		service.transactionManager.InjectTransaction(tx).(*server.TransactionManager))
}

func (service *AclService) OrganizationOwnerTemplate(organizatinId string) {
	actions := []string{"read", "delete", "modify", "create:folder"}

	service.CreatePolicyForResource("organization", organizatinId, actions)
}

func (service *AclService) CreatePolicyForResource(resourceName string, resourceId string, actions []string) (*Policy, error) {
	policy, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*AclService)

		policy := &Policy{}
		policy.Id = uuid.NewV4().String()

		if _, err := injectedService.aclRepository.InsertPolicy(policy); err != nil {
			return nil, err
		}

		for _, action := range actions {
			statement := &Statement{
				ResourceName: resourceName,
				ResourceID:   resourceId,
				Action:       action,
			}
			statement.Id = uuid.NewV4().String()

			if _, err := injectedService.aclRepository.InsertStatement(statement); err != nil {
				return nil, err
			}
			if err := injectedService.aclRepository.LinkStatementToPolicy(policy, statement); err != nil {
				return nil, err
			}
		}

		return policy, nil
	})

	return policy.(*Policy), err
}