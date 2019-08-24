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

func (repo *FolderService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewFolderService(repo.folderRepository.InjectTransaction(tx).(*FolderRepository),
		repo.organizationService.InjectTransaction(tx).(*OrganizationService),
		repo.aclService.InjectTransaction(tx).(*acl.AclService))
}

func (repo *FolderService) FindById(id string) *model.Folder {
	return repo.folderRepository.FindById(id)
}

func (repo *FolderService) Create(user *model.User, name string, organizationId string, parentFolderId *string) (*model.Folder, error) {

	// lets make sure the parent folder exists
	if parentFolderId != nil {
		parentFolder := repo.folderRepository.FindById(*parentFolderId)
		if parentFolder == nil {
			return nil, model.NewBadRequestError("could not find parent folder")
		}
	}

	folder := &model.Folder{
		Name:           name,
		OrganizationId: organizationId,
		ParentFolderId: parentFolderId,
	}
	folder.Id = uuid.NewV4().String()

	// check that they are allowed to create the folder in this organization
	canCreate, _ := repo.aclService.UserCanAccessResource(user, []string{"organization"}, []string{folder.OrganizationId}, "create:folder")
	if !canCreate {
		return nil, model.NewForbiddenError("you do not have permission to create a folder")
	}

	// we don't care about given this user specific access to this folder, they should keep the access because they have
	// it from the organization
	err := repo.folderRepository.Insert(folder)
	if err != nil {
		return nil, model.NewInternalServerError("failed to create folder")
	}

	return folder, nil
}

func (repo *FolderService) Update(user *model.User, folderId string, name string) (*model.Folder, error) {
	folder := repo.FindById(folderId)
	if folder == nil {
		return nil, model.NewNotFoundError("could not find folder")
	}

	resourceData, err := repo.aclService.GetResourceDataForModel(folder)
	if err != nil {
		return nil, model.NewInternalServerError("failed to fetch folder information")
	}

	// check that they are allowed to create the folder in this organization
	canUpdate, _ := repo.aclService.UserCanAccessResource(user, resourceData.ResourcePath, resourceData.ResourceIds, "modify")
	if !canUpdate {
		return nil, model.NewForbiddenError("you do not have permission to create a folder")
	}

	folder.Name = name

	err = repo.folderRepository.Update(folder)
	if err != nil {
		return nil, model.NewInternalServerError("failed to update folder")
	}

	return folder, nil
}

func (repo *FolderService) List(user *model.User, organizationId string, parentFolderId *string, pagination *model.Pagination) ([]model.Folder, error) {
	org := repo.organizationService.FindById(organizationId)
	if org == nil {
		return nil, model.NewNotFoundError("could not find organization")
	}

	orgResourceData, err := repo.aclService.GetResourceDataForModel(org)
	if err != nil {
		return nil, model.NewInternalServerError("failed to find organization information")
	}

	canAccess, err := repo.aclService.UserCanAccessResource(user, orgResourceData.ResourcePath, orgResourceData.ResourceIds, "view:folder")
	if err != nil {
		return nil, model.NewInternalServerError("could not verify access to organization")
	}

	if !canAccess {
		return nil, model.NewForbiddenError("you can not view folders in this organization")
	}

	folderResourceData, err := repo.aclService.GetResourceDataForModel(&model.Folder{})
	if err != nil {
		return nil, model.NewInternalServerError("failed to find folder information")
	}

	// find all of the resources that you can view
	resp, err := repo.aclService.UserActionableResourcesByPath(user, folderResourceData.ResourcePath, "view")
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

	folders, err := repo.folderRepository.Find(organizationIds, folderIds, parentFolderId, pagination)
	if err != nil {
		return nil, model.NewInternalServerError("failed to find folders")
	}

	return folders, nil
}