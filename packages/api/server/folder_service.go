package server

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/acl"
	"github.com/honerlaw/mentordoc/server/model"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type FolderService struct {
	folderRepository    *FolderRepository
	organizationService *OrganizationService
	aclService          *acl.AclService
}

func NewFolderService(folderRepository *FolderRepository, organizationService *OrganizationService, aclService *acl.AclService) *FolderService {
	return &FolderService{
		folderRepository:    folderRepository,
		organizationService: organizationService,
		aclService:          aclService,
	}
}

func (service *FolderService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewFolderService(service.folderRepository.InjectTransaction(tx).(*FolderRepository),
		service.organizationService.InjectTransaction(tx).(*OrganizationService),
		service.aclService.InjectTransaction(tx).(*acl.AclService))
}

func (service *FolderService) FindById(id string) *model.Folder {
	return service.folderRepository.FindById(id)
}

func (service *FolderService) Create(user *model.User, name string, organizationId string, parentFolderId *string) (*model.Folder, error) {

	// lets make sure the parent folder exists
	if parentFolderId != nil {
		parentFolder := service.folderRepository.FindById(*parentFolderId)
		if parentFolder == nil {
			return nil, model.NewBadRequestError("could not find parent folder")
		}
	}

	// make sure the org exists
	org := service.organizationService.FindById(organizationId)
	if org == nil {
		return nil, model.NewNotFoundError("could not find organization")
	}

	folder := &model.Folder{
		Name:           name,
		OrganizationId: org.Id,
		ParentFolderId: parentFolderId,
	}
	folder.Id = uuid.NewV4().String()

	// check that they are allowed to create the folder in this organization
	canCreate := service.aclService.UserCanAccessResourceByModel(user, org, "create:folder")
	if !canCreate {
		return nil, model.NewForbiddenError("you do not have permission to create a folder")
	}

	// we don't care about given this user specific access to this folder, they should keep the access because they have
	// it from the organization
	err := service.folderRepository.Insert(folder)
	if err != nil {
		return nil, model.NewInternalServerError("failed to create folder")
	}

	return folder, nil
}

func (service *FolderService) Update(user *model.User, folderId string, name string) (*model.Folder, error) {
	folder := service.FindById(folderId)
	if folder == nil {
		return nil, model.NewNotFoundError("could not find folder")
	}

	// check that they are allowed to create the folder in this organization
	canUpdate := service.aclService.UserCanAccessResourceByModel(user, folder, "modify")
	if !canUpdate {
		return nil, model.NewForbiddenError("you do not have permission to create a folder")
	}

	folder.Name = name

	err := service.folderRepository.Update(folder)
	if err != nil {
		return nil, model.NewInternalServerError("failed to update folder")
	}

	return folder, nil
}

func (service *FolderService) List(user *model.User, organizationId string, parentFolderId *string, pagination *model.Pagination) ([]model.Folder, error) {
	org := service.organizationService.FindById(organizationId)
	if org == nil {
		return nil, model.NewNotFoundError("could not find organization")
	}

	canAccess := service.aclService.UserCanAccessResourceByModel(user, org, "view:folder")
	if !canAccess {
		return nil, model.NewForbiddenError("you can not view folders in this organization")
	}

	folderResourceData, err := service.aclService.GetResourceDataForModel(&model.Folder{})
	if err != nil {
		return nil, model.NewInternalServerError("failed to find folder information")
	}

	// find all of the resources that you can view
	resp, err := service.aclService.UserActionableResourcesByPath(user, folderResourceData.ResourcePath, "view")
	if err != nil {
		return nil, model.NewInternalServerError("failed to find accessible folders")
	}

	organizationIds := make([]string, 0)
	folderIds := make([]string, 0)
	for _, res := range resp {
		if strings.HasPrefix(res.ResourcePath, "organization") {
			organizationIds = append(organizationIds, res.ResourceId)
		}
		if strings.HasPrefix(res.ResourcePath, "folder") {
			folderIds = append(folderIds, res.ResourceId)
		}
	}

	folders, err := service.folderRepository.Find(organizationIds, folderIds, parentFolderId, pagination)
	if err != nil {
		return nil, model.NewInternalServerError("failed to find folders")
	}

	return folders, nil
}