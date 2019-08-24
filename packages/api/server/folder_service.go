package server

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/acl"
	"github.com/honerlaw/mentordoc/server/model"
	uuid "github.com/satori/go.uuid"
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

func (repo *FolderService) Create(user *model.User, name string, organizationId string, parentFolderId *string) (*model.Folder, error) {

	// lets make sure the parent folder exists
	if parentFolderId != nil {
		parentFolder := repo.folderRepository.FindById(*parentFolderId)
		if parentFolder == nil {
			return nil, model.NewBadRequestError("could not find parent folder")
		}
	}

	// check that they are allowed to create the folder in this organization
	canCreate, _ := repo.aclService.UserCanAccessResource(user, []string{"organization"}, []string{organizationId}, "create:folder")
	if !canCreate {
		return nil, model.NewForbiddenError("you do not have permission to create a folder")
	}

	folder := &model.Folder{
		Name:           name,
		OrganizationId: organizationId,
		ParentFolderId: parentFolderId,
	}
	folder.Id = uuid.NewV4().String()

	// we don't care about given this user specific access to this folder, they should keep the access because they have
	// it from the organization
	err := repo.folderRepository.Insert(folder)
	if err != nil {
		return nil, model.NewInternalServerError("failed to create folder")
	}

	return folder, nil
}
