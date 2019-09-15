package document

import (
	"database/sql"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/folder"
	"github.com/honerlaw/mentordoc/server/lib/organization"
	"github.com/honerlaw/mentordoc/server/lib/resource_history"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type DocumentService struct {
	documentRepository        *DocumentRepository
	documentDraftRepository   *DocumentDraftRepository
	documentContentRepository *DocumentContentRepository
	organizationService       *organization.OrganizationService
	folderService             *folder.FolderService
	aclService                *acl.AclService
	transactionManager        *util.TransactionManager
	resourceHistoryService    *resource_history.ResourceHistoryService
}

func NewDocumentService(
	documentRepository *DocumentRepository,
	documentDraftRepository *DocumentDraftRepository,
	documentContentRepository *DocumentContentRepository,
	organizationService *organization.OrganizationService,
	folderService *folder.FolderService,
	aclService *acl.AclService,
	transactionManager *util.TransactionManager,
	resourceHistoryService *resource_history.ResourceHistoryService,
) *DocumentService {
	return &DocumentService{
		documentRepository:        documentRepository,
		documentDraftRepository:   documentDraftRepository,
		documentContentRepository: documentContentRepository,
		organizationService:       organizationService,
		folderService:             folderService,
		aclService:                aclService,
		transactionManager:        transactionManager,
		resourceHistoryService:    resourceHistoryService,
	}
}

func (service *DocumentService) InjectTransaction(tx *sql.Tx) interface{} {
	return NewDocumentService(
		service.documentRepository.InjectTransaction(tx).(*DocumentRepository),
		service.documentDraftRepository.InjectTransaction(tx).(*DocumentDraftRepository),
		service.documentContentRepository.InjectTransaction(tx).(*DocumentContentRepository),
		service.organizationService.InjectTransaction(tx).(*organization.OrganizationService),
		service.folderService.InjectTransaction(tx).(*folder.FolderService),
		service.aclService.InjectTransaction(tx).(*acl.AclService),
		service.transactionManager.InjectTransaction(tx).(*util.TransactionManager),
		service.resourceHistoryService.InjectTransaction(tx).(*resource_history.ResourceHistoryService),
	)
}

func (service *DocumentService) FindDocument(user *shared.User, documentId string) (*shared.Document, error) {
	document := service.documentRepository.FindById(documentId)
	if document == nil {
		return nil, shared.NewNotFoundError("could not find document")
	}

	// document has not been initially published, so find the draft for this user
	if document.InitialDraftUserId != nil {
		if *document.InitialDraftUserId != user.Id {
			return nil, shared.NewForbiddenError("can not access document")
		}
		return service.FindDraftDocument(user, document.Id)
	}

	canAccess := service.aclService.UserCanAccessResourceByModel(user, document, "view", "modify")
	if !canAccess {
		return nil, shared.NewForbiddenError("can not view document")
	}

	draft := service.documentDraftRepository.FindPublishedDraftByDocumentId(documentId)
	if draft == nil {
		return nil, shared.NewNotFoundError("could not find published document")
	}

	content := service.documentContentRepository.FindByDocumentDraftId(draft.Id)
	if content == nil {
		return nil, shared.NewNotFoundError("could not find document content");
	}

	draft.Content = content
	document.Drafts = []shared.DocumentDraft{*draft}

	return document, nil
}

func (service *DocumentService) FindDraftDocument(user *shared.User, documentId string) (*shared.Document, error) {
	document := service.documentRepository.FindById(documentId)
	if document == nil {
		return nil, shared.NewNotFoundError("could not find document")
	}

	canAccess := service.aclService.UserCanAccessResourceByModel(user, document, "modify")
	if !canAccess {
		return nil, shared.NewForbiddenError("can not modify document")
	}

	draft := service.documentDraftRepository.FindDraftByDocumentId(documentId)
	if draft == nil {
		return nil, shared.NewNotFoundError("could not find published document")
	}

	content := service.documentContentRepository.FindByDocumentDraftId(draft.Id)
	if content == nil {
		return nil, shared.NewNotFoundError("could not find document content");
	}

	draft.Content = content
	document.Drafts = []shared.DocumentDraft{*draft}

	return document, nil
}

func (service *DocumentService) FindDocumentAncestry(user *shared.User, documentId string) ([]interface{}, error) {
	document, err := service.FindDocument(user, documentId)
	if err != nil {
		return nil, err
	}

	path := make([]interface{}, 0)
	path = append(path, document)
	if document.FolderId != nil {
		folders, err := service.folderService.FindAncestry(*document.FolderId)
		if err != nil {
			return nil, shared.NewInternalServerError("failed to find document path");
		}
		for _, f := range folders {
			path = append(path, f)
		}
	}

	org := service.organizationService.FindById(document.OrganizationId)
	if org == nil {
		return nil, shared.NewNotFoundError("could not find organization")
	}

	path = append(path, org)

	// reverse the array
	for i := len(path) / 2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}

	return path, nil
}

func (service *DocumentService) Create(user *shared.User, organizationId string, folderId *string, name string, content string) (*shared.Document, error) {
	organizationId, folderId, err := service.hasAccessToOrganizationOrFolder(user, organizationId, folderId, "create:document")
	if err != nil {
		return nil, err
	}

	document := &shared.Document{
		OrganizationId:     organizationId,
		FolderId:           folderId,
		InitialDraftUserId: &user.Id,
	}
	document.Id = uuid.NewV4().String()

	documentDraft := &shared.DocumentDraft{
		DocumentId: document.Id,
		Name:       name,
	}
	documentDraft.Id = uuid.NewV4().String()

	documentContent := &shared.DocumentContent{
		DocumentDraftId: documentDraft.Id,
		Content:         content,
	}
	documentContent.Id = uuid.NewV4().String()

	_, err = service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*DocumentService)

		err := injectedService.documentRepository.Insert(document)
		if err != nil {
			return nil, err
		}

		err = injectedService.documentDraftRepository.Insert(documentDraft)
		if err != nil {
			return nil, err
		}

		err = injectedService.documentContentRepository.Insert(documentContent)
		if err != nil {
			return nil, err
		}

		_, err = injectedService.resourceHistoryService.Create(document.Id, "document", user.Id, "created")
		if err != nil {
			return nil, err
		}

		_, err = injectedService.resourceHistoryService.Create(documentDraft.Id, "document_draft", user.Id, "created")
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return nil, shared.NewInternalServerError("failed to create document")
	}

	documentDraft.Content = documentContent
	document.Drafts = []shared.DocumentDraft{*documentDraft}

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

	// find the draft version, we can only update if a draft exists
	documentDraft := service.documentDraftRepository.FindDraftByDocumentId(document.Id)
	if documentDraft == nil {
		return nil, shared.NewBadRequestError("could not find draft version off the document to update")
	}

	documentContent := service.documentContentRepository.FindByDocumentDraftId(documentDraft.Id)
	if documentContent == nil {
		return nil, shared.NewNotFoundError("could not find document content")
	}

	res, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*DocumentService)

		documentDraft.Name = name
		err := injectedService.documentDraftRepository.Update(documentDraft)
		if err != nil {
			return nil, err
		}

		documentContent.Content = content
		err = injectedService.documentContentRepository.Update(documentContent)
		if err != nil {
			return nil, err
		}

		_, err = injectedService.resourceHistoryService.Create(document.Id, "document", user.Id, "updated")
		if err != nil {
			return nil, err
		}

		_, err = injectedService.resourceHistoryService.Create(documentDraft.Id, "document_draft", user.Id, "updated")
		if err != nil {
			return nil, err
		}

		documentDraft.Content = documentContent
		document.Drafts = []shared.DocumentDraft{*documentDraft}

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

	res, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*DocumentService)

		deletedAt := util.NowUnix()

		document.DeletedAt = &deletedAt
		err := injectedService.documentRepository.Update(document)
		if err != nil {
			return nil, err
		}

		// delete all of the drafts as well, we implicitly delete the content for each draft this way
		err = injectedService.documentDraftRepository.Delete(document.Id);
		if err != nil {
			return nil, err
		}

		_, err = injectedService.resourceHistoryService.Create(document.Id, "document", user.Id, "deleted")
		if err != nil {
			return nil, err
		}

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
	resp, err := service.aclService.UserActionableResourcesByPath(user, documentResourceData.ResourcePath, "view", "modify")
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

	documents, err := service.documentRepository.Find(user.Id, organizationIds, folderIds, documentIds, folderId, pagination)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find documents")
	}

	err = service.documentDraftRepository.FindAndAttachLatestDraftForDocuments(documents)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find document drafts")
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
		fold := service.folderService.FindById(*folderId)
		if fold == nil {
			return "", nil, shared.NewNotFoundError("could not find folder")
		}

		canAccess := service.aclService.UserCanAccessResourceByModel(user, fold, action)
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
