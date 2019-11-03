package controller

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/http/middleware"
	"github.com/honerlaw/mentordoc/server/http/request"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/document"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"net/http"
)

type DocumentController struct {
	validatorService         *util.ValidatorService
	documentService          *document.DocumentService
	authenticationMiddleware *middleware.AuthenticationMiddleware
	aclService               *acl.AclService
}

func NewDocumentController(
	validatorService *util.ValidatorService,
	documentService *document.DocumentService,
	authenticationMiddleware *middleware.AuthenticationMiddleware,
	aclService *acl.AclService,
) *DocumentController {

	return &DocumentController{
		validatorService:         validatorService,
		documentService:          documentService,
		authenticationMiddleware: authenticationMiddleware,
		aclService:               aclService,
	}
}

func (controller *DocumentController) RegisterRoutes(router chi.Router) {
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/document/{id}", controller.find)
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/document/path/{id}", controller.findPath)
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/document/search", controller.search)
	router.
		With(controller.validatorService.Middleware(request.DocumentCreateRequest{}), controller.authenticationMiddleware.HasAccessToken()).
		Post("/document", controller.create)
	router.
		With(controller.validatorService.Middleware(request.DocumentDraftCreateRequest{}), controller.authenticationMiddleware.HasAccessToken()).
		Post("/document/draft", controller.createDraft)
	router.
		With(controller.validatorService.Middleware(request.DocumentUpdateRequest{}), controller.authenticationMiddleware.HasAccessToken()).
		Put("/document", controller.update)
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/document/list/{organizationId}", controller.list)
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Delete("/document/{id}", controller.delete)
}

func (controller *DocumentController) find(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	documentId := chi.URLParam(req, "id")

	doc, err := controller.documentService.FindDocument(user, documentId)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, []*shared.Document{doc})
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("found document but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusCreated, wrapped[0])
}

func (controller *DocumentController) findPath(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	documentId := chi.URLParam(req, "id")

	path, err := controller.documentService.FindDocumentAncestry(user, documentId)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, path)
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("could not find user access for items in path"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusCreated, wrapped)
}

func (controller *DocumentController) create(w http.ResponseWriter, req *http.Request) {
	validReq := controller.validatorService.GetModelFromRequest(req).(*request.DocumentCreateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)

	doc, err := controller.documentService.Create(user, validReq.OrganizationId, validReq.FolderId, validReq.Name, validReq.Content)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, []*shared.Document{doc})
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("created document but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusCreated, wrapped[0])
}

func (controller *DocumentController) createDraft(w http.ResponseWriter, req *http.Request) {
	validReq := controller.validatorService.GetModelFromRequest(req).(*request.DocumentDraftCreateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)

	doc, err := controller.documentService.CreateDraft(user, validReq.DocumentId, validReq.Name, validReq.Content)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, []*shared.Document{doc})
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("created document draft but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusCreated, wrapped[0])
}

func (controller *DocumentController) update(w http.ResponseWriter, req *http.Request) {
	validReq := controller.validatorService.GetModelFromRequest(req).(*request.DocumentUpdateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)

	doc, err := controller.documentService.Update(user, validReq.DocumentId, validReq.DraftId,
		validReq.Name, validReq.Content, validReq.ShouldPublish, validReq.ShouldRetract)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, []*shared.Document{doc})
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("updated document but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped[0])
}

func (controller *DocumentController) delete(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	documentId := chi.URLParam(req, "id")

	doc, err := controller.documentService.Delete(user, documentId)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, []*shared.Document{doc})
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("updated document but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped[0])
}

func (controller *DocumentController) list(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	organizationId := chi.URLParam(req, "organizationId")
	queryFolderId := req.URL.Query().Get("folderId")
	pagination := shared.NewPagination(req)

	if len(organizationId) == 0 {
		util.WriteHttpError(w, shared.NewBadRequestError("organization is required"))
		return
	}

	var folderId *string
	if len(queryFolderId) > 0 {
		folderId = &queryFolderId
	}

	documents, err := controller.documentService.List(user, organizationId, folderId, pagination)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	if len(documents) == 0 {
		util.WriteJsonToResponse(w, http.StatusOK, documents)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, documents)
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("found documents but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped)
}

func (controller *DocumentController) search(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	searchQuery := req.URL.Query().Get("query")

	if len(searchQuery) == 0 {
		util.WriteHttpError(w, shared.NewBadRequestError("a search query is required"))
		return
	}

	documents, err := controller.documentService.Search(user, searchQuery);
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("failed to find documents"))
		return
	}

	wrapped, err := controller.aclService.Wrap(user, documents)
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("found documents but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped)
}