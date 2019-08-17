package server

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/transaction"
	uuid "github.com/satori/go.uuid"
)

type OrganizationService struct {
	transaction.Transactionable
	organizationRepository *OrganizationRepository
}

func NewOrganizationService(organizationRepository *OrganizationRepository) *OrganizationService {
	service := &OrganizationService{
		organizationRepository: organizationRepository,
	};
	service.CloneWithTransaction = func(tx *sql.Tx) interface{} {
		return NewOrganizationService(service.organizationRepository.InjectTransaction(tx).(*OrganizationRepository))
	}
	return service
}

func (service *OrganizationService) Create(name string) (*Organization, error) {
	organization := &Organization{
		Name: name,
	}
	organization.Id = uuid.NewV4().String()

	return service.organizationRepository.Insert(organization)
}
