package server

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/acl"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	"net/http"
)

type DocumentController struct {
	validatorService         *util.ValidatorService
	documentService          *DocumentService
	authenticationMiddleware *AuthenticationMiddleware
	aclService               *acl.AclService
}

func NewDocumentController(
	validatorService *util.ValidatorService,
	documentService *DocumentService,
	authenticationMiddleware *AuthenticationMiddleware,
	aclService *acl.AclService,
) *DocumentController {

	return &DocumentController{
		validatorService:         validatorService,
		documentService:          documentService,
		authenticationMiddleware: authenticationMiddleware,
		aclService:               aclService,
	}
}

/**

Featuers
- Create
- Update
- List

Future
- Verification
- Drafts
*/
func (controller *DocumentController) RegisterRoutes(router chi.Router) {
	router.
		With(controller.validatorService.Middleware(model.DocumentCreateRequest{}), controller.authenticationMiddleware.HasAccessToken()).
		Post("/document", controller.create)
	router.
		With(controller.validatorService.Middleware(model.DocumentUpdateRequest{}), controller.authenticationMiddleware.HasAccessToken()).
		Put("/document", controller.update)
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/document/list/{organizationId}", controller.list)
}

func (controller *DocumentController) create(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*model.DocumentCreateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)


	document, err := controller.documentService.Create(user, request.OrganizationId, request.FolderId, request.Name, request.Content)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, []*model.Document{document})
	if err != nil {
		util.WriteHttpError(w, model.NewInternalServerError("created document but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusCreated, wrapped[0])
}

func (controller *DocumentController) update(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*model.DocumentUpdateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)

	document, err := controller.documentService.Update(user, request.DocumentId, request.Name, request.Content)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, []*model.Document{document})
	if err != nil {
		util.WriteHttpError(w, model.NewInternalServerError("updated document but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped[0])
}

func (controller *DocumentController) list(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	organizationId := chi.URLParam(req, "organizationId")
	queryFolderId := req.URL.Query().Get("folderId")
	pagination := model.NewPagination(req)

	if len(organizationId) == 0 {
		util.WriteHttpError(w, model.NewBadRequestError("organization is required"))
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
		util.WriteHttpError(w, model.NewInternalServerError("found documents but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped)
}
