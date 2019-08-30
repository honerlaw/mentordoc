package organization

import (
	"database/sql"
	uuid "github.com/satori/go.uuid"
)

type OrganizationService struct {
	organizationRepository *OrganizationRepository
}

func NewOrganizationService(organizationRepository *OrganizationRepository) *OrganizationService {
	service := &OrganizationService{
		organizationRepository: organizationRepository,
	};
	return service
}

func (service *OrganizationService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewOrganizationService(service.organizationRepository.InjectTransaction(tx).(*OrganizationRepository))
}

func (service *OrganizationService) Create(name string) (*Organization, error) {
	organization := &Organization{
		Name: name,
	}
	organization.Id = uuid.NewV4().String()

	return service.organizationRepository.Insert(organization)
}

func (service *OrganizationService) FindById(id string) *Organization {
	return service.organizationRepository.FindById(id)
}
