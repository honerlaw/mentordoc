package organization

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type OrganizationService struct {
	organizationRepository *OrganizationRepository
	aclService             *acl.AclService
}

func NewOrganizationService(
	organizationRepository *OrganizationRepository,
	aclService *acl.AclService,
) *OrganizationService {
	service := &OrganizationService{
		organizationRepository: organizationRepository,
		aclService:             aclService,
	};
	return service
}

func (service *OrganizationService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewOrganizationService(
		service.organizationRepository.InjectTransaction(tx).(*OrganizationRepository),
		service.aclService.InjectTransaction(tx).(*acl.AclService),
	)
}

func (service *OrganizationService) Create(name string) (*shared.Organization, error) {
	organization := &shared.Organization{
		Name: name,
	}
	organization.Id = uuid.NewV4().String()

	return service.organizationRepository.Insert(organization)
}

func (service *OrganizationService) List(u *shared.User) ([]shared.Organization, error) {

	orgResourceData, err := service.aclService.GetResourceDataForModel(&shared.Organization{})
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find organization information")
	}

	// find all of the resources that you can view
	resp, err := service.aclService.UserActionableResourcesByPath(u, orgResourceData.ResourcePath, "view")
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find accessible organizations")
	}

	organizationIds := make([]string, 0)
	for _, res := range resp {
		if strings.HasPrefix(res.ResourcePath, "organization") {
			organizationIds = append(organizationIds, res.ResourceId)
		}
	}

	orgs, err := service.organizationRepository.Find(organizationIds)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find organizations")
	}

	return orgs, nil
}

func (service *OrganizationService) FindById(id string) *shared.Organization {
	return service.organizationRepository.FindById(id)
}
