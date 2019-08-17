package server

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
)

type OrganizationService struct {
	Transactionable
	organizationRepository *OrganizationRepository
}

func NewOrganizationService(organizationRepository *OrganizationRepository) *OrganizationService {
	service := &OrganizationService{
		organizationRepository: organizationRepository,
	};
	service.cloneWithTransaction = func(tx *sql.Tx) interface{} {
		return NewOrganizationService(service.organizationRepository.InjectTransaction(tx).(*OrganizationRepository))
	}
	return service
}

func (service *OrganizationService) Create(name string) (*Organization, error) {
	organization := &Organization{
		Name: name,
	}
	organization.Id = uuid.NewV4().String()

	return nil, nil
}
