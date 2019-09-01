package folder

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/organization"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type FolderService struct {
	folderRepository    *FolderRepository
	organizationService *organization.OrganizationService
	aclService          *acl.AclService
}

func NewFolderService(folderRepository *FolderRepository, organizationService *organization.OrganizationService, aclService *acl.AclService) *FolderService {
	return &FolderService{
		folderRepository:    folderRepository,
		organizationService: organizationService,
		aclService:          aclService,
	}
}

func (service *FolderService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewFolderService(service.folderRepository.InjectTransaction(tx).(*FolderRepository),
		service.organizationService.InjectTransaction(tx).(*organization.OrganizationService),
		service.aclService.InjectTransaction(tx).(*acl.AclService))
}

func (service *FolderService) FindById(id string) *shared.Folder {
	return service.folderRepository.FindById(id)
}

func (service *FolderService) Create(user *shared.User, name string, organizationId string, parentFolderId *string) (*shared.Folder, error) {

	// lets make sure the parent folder exists
	if parentFolderId != nil {
		parentFolder := service.folderRepository.FindById(*parentFolderId)
		if parentFolder == nil {
			return nil, shared.NewBadRequestError("could not find parent folder")
		}
	}

	// make sure the org exists
	org := service.organizationService.FindById(organizationId)
	if org == nil {
		return nil, shared.NewNotFoundError("could not find organization")
	}

	folder := &shared.Folder{
		Name:           name,
		OrganizationId: org.Id,
		ParentFolderId: parentFolderId,
	}
	folder.Id = uuid.NewV4().String()

	// check that they are allowed to create the folder in this organization
	canCreate := service.aclService.UserCanAccessResourceByModel(user, org, "create:folder")
	if !canCreate {
		return nil, shared.NewForbiddenError("you do not have permission to create a folder")
	}

	// we don't care about given this user specific access to this folder, they should keep the access because they have
	// it from the organization
	err := service.folderRepository.Insert(folder)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to create folder")
	}

	return folder, nil
}

func (service *FolderService) Update(user *shared.User, folderId string, name string) (*shared.Folder, error) {
	folder := service.FindById(folderId)
	if folder == nil {
		return nil, shared.NewNotFoundError("could not find folder")
	}

	// check that they are allowed to create the folder in this organization
	canUpdate := service.aclService.UserCanAccessResourceByModel(user, folder, "modify")
	if !canUpdate {
		return nil, shared.NewForbiddenError("you do not have permission to create a folder")
	}

	folder.Name = name

	err := service.folderRepository.Update(folder)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to update folder")
	}

	return folder, nil
}

func (service *FolderService) Delete(user *shared.User, folderId string) (*shared.Folder, error) {
	folder := service.FindById(folderId)
	if folder == nil {
		return nil, shared.NewNotFoundError("could not find folder")
	}

	if folder.ChildCount > 0 {
		return nil, shared.NewBadRequestError("can not delete a folder with contents")
	}

	// check that they are allowed to create the folder in this organization
	canDelete := service.aclService.UserCanAccessResourceByModel(user, folder, "delete")
	if !canDelete {
		return nil, shared.NewForbiddenError("you do not have permission to delete this folder")
	}

	deletedAt := util.NowUnix()
	folder.DeletedAt = &deletedAt

	err := service.folderRepository.Update(folder)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to delete folder")
	}

	return folder, nil
}


func (service *FolderService) List(user *shared.User, organizationId string, parentFolderId *string, pagination *shared.Pagination) ([]shared.Folder, error) {
	org := service.organizationService.FindById(organizationId)
	if org == nil {
		return nil, shared.NewNotFoundError("could not find organization")
	}

	canAccess := service.aclService.UserCanAccessResourceByModel(user, org, "view:folder")
	if !canAccess {
		return nil, shared.NewForbiddenError("you can not view folders in this organization")
	}

	folderResourceData, err := service.aclService.GetResourceDataForModel(&shared.Folder{})
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find folder information")
	}

	// find all of the resources that you can view
	resp, err := service.aclService.UserActionableResourcesByPath(user, folderResourceData.ResourcePath, "view")
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find accessible folders")
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
		return nil, shared.NewInternalServerError("failed to find folders")
	}

	return folders, nil
}