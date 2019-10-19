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

	// make sure the document itself is accessible
	canAccess := service.aclService.UserCanAccessResourceByModel(user, document, "view", "modify")
	if !canAccess {
		return nil, shared.NewForbiddenError("can not view document")
	}

	// find the latest draft that we can access / view
	drafts, err := service.documentDraftRepository.FindLatestAccessibleDraftForDocuments(user.Id, []string{document.Id});

	// no drafts were found that this user can access
	if err != nil || len(drafts) == 0 {
		return nil, shared.NewForbiddenError("can not access document")
	}

	// there should only ever be one draft in this scenario (the latest accessible one)
	draft := &drafts[0]

	// attach the content to the draft
	content := service.documentContentRepository.FindByDocumentDraftId(draft.Id)
	if content == nil {
		return nil, shared.NewNotFoundError("could not find document content");
	}

	draft.Content = content
	document.Drafts = drafts

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
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}

	return path, nil
}

func (service *DocumentService) CreateDraft(user *shared.User, documentId string, name string, content string) (*shared.Document, error) {
	document := service.documentRepository.FindById(documentId)
	if document == nil {
		return nil, shared.NewNotFoundError("could not find document")
	}

	canAccess := service.aclService.UserCanAccessResourceByModel(user, document, "modify")
	if !canAccess {
		return nil, shared.NewForbiddenError("can not modify document")
	}

	documentDraft := &shared.DocumentDraft{
		DocumentId: document.Id,
		Name:       name,
		CreatorId:  user.Id,
	}
	documentDraft.Id = uuid.NewV4().String()

	documentContent := &shared.DocumentContent{
		DocumentDraftId: documentDraft.Id,
		Content:         content,
	}
	documentContent.Id = uuid.NewV4().String()

	_, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*DocumentService)

		err := injectedService.documentDraftRepository.Insert(documentDraft)
		if err != nil {
			return nil, err
		}

		err = injectedService.documentContentRepository.Insert(documentContent)
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
		return nil, shared.NewInternalServerError("failed to create document draft")
	}

	documentDraft.Content = documentContent
	document.Drafts = []shared.DocumentDraft{*documentDraft}

	return document, nil
}

func (service *DocumentService) Create(user *shared.User, organizationId string, folderId *string, name string, content string) (*shared.Document, error) {
	organizationId, folderId, err := service.hasAccessToOrganizationOrFolder(user, organizationId, folderId, "create:document")
	if err != nil {
		return nil, err
	}

	document := &shared.Document{
		OrganizationId: organizationId,
		FolderId:       folderId,
	}
	document.Id = uuid.NewV4().String()

	documentDraft := &shared.DocumentDraft{
		DocumentId: document.Id,
		Name:       name,
		CreatorId:  user.Id,
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

func (service *DocumentService) Update(
	user *shared.User, documentId string, draftId string,
	name *string, content *string, shouldPublish bool, shouldRetract bool,
) (*shared.Document, error) {
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

	if documentDraft.Id != draftId {
		return nil, shared.NewBadRequestError("target draft and current draft are not the same")
	}

	documentContent := service.documentContentRepository.FindByDocumentDraftId(documentDraft.Id)
	if documentContent == nil {
		return nil, shared.NewNotFoundError("could not find document content")
	}

	res, err := service.transactionManager.Transact(service, func(injected interface{}) (interface{}, error) {
		injectedService := injected.(*DocumentService)

		if name != nil {
			documentDraft.Name = *name
		}
		if shouldPublish {
			publishedAt := util.NowUnix()
			documentDraft.PublishedAt = &publishedAt
		}
		if shouldRetract {
			retractedAt := util.NowUnix()
			documentDraft.RetractedAt = &retractedAt
			documentDraft.DeletedAt = &retractedAt
		}
		err := injectedService.documentDraftRepository.Update(documentDraft)
		if err != nil {
			return nil, err
		}

		if content != nil {
			documentContent.Content = *content
		}
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

	// this will find all of the documents that you are able to view, but does not take into account drafts that are not tied to you
	documents, err := service.documentRepository.Find(user.Id, organizationIds, folderIds, documentIds, folderId, pagination)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find documents")
	}

	// so next up find the latest draft for each document, and attach them to the document
	err = service.documentDraftRepository.FindAndAttachLatestAccessibleDraftForDocuments(user.Id, documents)
	if err != nil {
		return nil, shared.NewInternalServerError("failed to find document drafts")
	}

	// make sure the doc has at least one valid draft on it to show
	validDocuments := make([]shared.Document, 0);
	for i := 0; i < len(documents); i++ {
		doc := documents[i]
		if len(doc.Drafts) > 0 {
			validDocuments = append(validDocuments, doc)
		}
	}

	return validDocuments, nil
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
