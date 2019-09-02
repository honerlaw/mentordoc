package document

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/folder"
	"github.com/honerlaw/mentordoc/server/lib/organization"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type DocumentService struct {
	documentRepository        *DocumentRepository
	documentContentRepository *DocumentContentRepository
	organizationService       *organization.OrganizationService
	folderService             *folder.FolderService
	aclService                *acl.AclService
	transactionManager        *util.TransactionManager
}

func NewDocumentService(
	documentRepository *DocumentRepository,
	documentContentRepository *DocumentContentRepository,
	organizationService *organization.OrganizationService,
	folderService *folder.FolderService,
	aclService *acl.AclService,
	transactionManager *util.TransactionManager,
) *DocumentService {
	return &DocumentService{
		documentRepository:        documentRepository,
		documentContentRepository: documentContentRepository,
		organizationService:       organizationService,
		folderService:             folderService,
		aclService:                aclService,
		transactionManager:        transactionManager,
	}
}

func (service *DocumentService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewDocumentService(
		service.documentRepository.InjectTransaction(tx).(*DocumentRepository),
		service.documentContentRepository.InjectTransaction(tx).(*DocumentContentRepository),
		service.organizationService.InjectTransaction(tx).(*organization.OrganizationService),
		service.folderService.InjectTransaction(tx).(*folder.FolderService),
		service.aclService.InjectTransaction(tx).(*acl.AclService),
		service.transactionManager.InjectTransaction(tx).(*util.TransactionManager),
	)
}

func (service *DocumentService) Find(user *shared.User, documentId string) (*shared.Document, error) {
	document := service.documentRepository.FindById(documentId)
	if document == nil {
		return nil, shared.NewNotFoundError("could not find document")
	}

	canAccess := service.aclService.UserCanAccessResourceByModel(user, document, "view")
	if !canAccess {
		return nil, shared.NewForbiddenError("can not modify document")
	}

	content := service.documentContentRepository.FindByDocumentId(documentId)
	if content == nil {
		return nil, shared.NewNotFoundError("could not find document content");
	}

	document.Content = content

	return document, nil
}

func (service *DocumentService) Create(user *shared.User, organizationId string, folderId *string, name string, content string) (*shared.Document, error) {
	organizationId, folderId, err := service.hasAccessToOrganizationOrFolder(user, organizationId, folderId, "create:document")
	if err != nil {
		return nil, err
	}

	document := &shared.Document{
		Name:           name,
		OrganizationId: organizationId,
		FolderId:       folderId,
	}
	document.Id = uuid.NewV4().String()

	documentContent := &shared.DocumentContent{
		DocumentId: document.Id,
		Content:    content,
	}
	documentContent.Id = uuid.NewV4().String()

	_, err = service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*DocumentService)

		err := injectedService.documentRepository.Insert(document)
		if err != nil {
			return nil, err
		}

		err = injectedService.documentContentRepository.Insert(documentContent)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return nil, shared.NewInternalServerError("failed to create document")
	}

	document.Content = documentContent

	return document, nil
}

func (service *DocumentService) Update(user *shared.User, documentId string, name string, content string) (*shared.Document, error) {
	document := service.documentRepository.FindById(documentId)
	if document == nil {
		return nil, shared.NewNotFoundError("could not find document")
	}

	canAccess := service.aclService.UserCanAccessResourceByModel(user, document, "modify")
	if !canAccess {
		return nil, shared.NewForbiddenError("can not modify document")
	}

	documentContent := service.documentContentRepository.FindByDocumentId(document.Id)
	if documentContent == nil {
		return nil, shared.NewNotFoundError("could not find document content")
	}

	res, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*DocumentService)

		document.Name = name
		err := injectedService.documentRepository.Update(document)
		if err != nil {
			return nil, err
		}

		documentContent.Content = content
		err = injectedService.documentContentRepository.Update(documentContent)
		if err != nil {
			return nil, err
		}

		document.Content = documentContent

		return document, nil
	})

	if err != nil {
		return nil, shared.NewInternalServerError("failed to update document")
	}

	return res.(*shared.Document), nil
}

func (service *DocumentService) Delete(user *shared.User, documentId string) (*shared.Document, error) {
	document := service.documentRepository.FindById(documentId)
	if document == nil {
		return nil, shared.NewNotFoundError("could not find document")
	}

	canAccess := service.aclService.UserCanAccessResourceByModel(user, document, "delete")
	if !canAccess {
		return nil, shared.NewForbiddenError("can not delete document")
	}

	documentContent := service.documentContentRepository.FindByDocumentId(document.Id)
	if documentContent == nil {
		return nil, shared.NewNotFoundError("could not find document content")
	}

	res, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*DocumentService)

		deletedAt := util.NowUnix()

		document.DeletedAt = &deletedAt
		err := injectedService.documentRepository.Update(document)
		if err != nil {
			return nil, err
		}

		documentContent.DeletedAt = &deletedAt
		err = injectedService.documentContentRepository.Update(documentContent)
		if err != nil {
			return nil, err
		}

		document.Content = documentContent

		return document, nil
	})

	if err != nil {
		return nil, shared.NewInternalServerError("failed to delete document")
	}

	return res.(*shared.Document), nil
}

func (service *DocumentService) List(user *shared.User, organizationId string, folderId *string, pagination *shared.Pagination) ([]shared.Document, error) {
	organizationId, folderId, err := service.hasAccessToOrganizationOrFolder(user, organizationId, folderId, "view:document");
	if err != nil {
		return nil, err
	}

	documentResourceData, err := service.aclService.GetResourceDataForModel(&shared.Document{})
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find document information")
	}

	// find all of the resources that you can view
	resp, err := service.aclService.UserActionableResourcesByPath(user, documentResourceData.ResourcePath, "view")
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find accessible documents")
	}

	organizationIds := make([]string, 0)
	folderIds := make([]string, 0)
	documentIds := make([]string, 0)
	for _, res := range resp {
		if strings.HasPrefix(res.ResourcePath, "organization") {
			organizationIds = append(organizationIds, res.ResourceId)
		}
		if strings.HasPrefix(res.ResourcePath, "folder") {
			folderIds = append(folderIds, res.ResourceId)
		}
		if strings.HasPrefix(res.ResourcePath, "document") {
			documentIds = append(documentIds, res.ResourceId)
		}
	}

	documents, err := service.documentRepository.Find(organizationIds, folderIds, documentIds, folderId, pagination)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find folders")
	}

	return documents, nil
}

func (service *DocumentService) hasAccessToOrganizationOrFolder(user *shared.User, organizationId string, folderId *string, action string) (string, *string, error) {
	org := service.organizationService.FindById(organizationId)
	if org == nil {
		return "", nil, shared.NewNotFoundError("could not find organization")
	}

	// if they are adding this to a folder, check the folder exists and they have access
	if folderId != nil {
		folder := service.folderService.FindById(*folderId)
		if folder == nil {
			return "", nil, shared.NewNotFoundError("could not find folder")
		}

		canAccess := service.aclService.UserCanAccessResourceByModel(user, folder, action)
		if !canAccess {
			return "", nil, shared.NewForbiddenError("can not create document in folder")
		}
	} else {
		// they are not trying to add this to a folder, make sure they can do it at the org level
		canAccess := service.aclService.UserCanAccessResourceByModel(user, org, action)
		if !canAccess {
			return "", nil, shared.NewForbiddenError("can not create document in organization")
		}
	}
	return organizationId, folderId, nil
}
