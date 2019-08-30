package controller

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/http/middleware"
	"github.com/honerlaw/mentordoc/server/http/request"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/folder"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"net/http"
)

type FolderController struct {
	validatorService         *util.ValidatorService
	folderService            *folder.FolderService
	authenticationMiddleware *middleware.AuthenticationMiddleware
	aclService               *acl.AclService
}

func NewFolderController(
	validatorService *util.ValidatorService,
	folderService *folder.FolderService,
	authenticationMiddleware *middleware.AuthenticationMiddleware,
	aclService *acl.AclService,
) *FolderController {

	return &FolderController{
		validatorService:         validatorService,
		folderService:            folderService,
		authenticationMiddleware: authenticationMiddleware,
		aclService:               aclService,
	}
}

func (controller *FolderController) RegisterRoutes(router chi.Router) {
	router.
		With(controller.validatorService.Middleware(request.FolderCreateRequest{}),
			controller.authenticationMiddleware.HasAccessToken()).
		Post("/folder", controller.create)
	router.
		With(controller.validatorService.Middleware(request.FolderUpdateRequest{}),
			controller.authenticationMiddleware.HasAccessToken()).
		Put("/folder/{id}", controller.update)
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/folder/list/{organizationId}", controller.list)
}

func (controller *FolderController) create(w http.ResponseWriter, req *http.Request) {
	validReq := controller.validatorService.GetModelFromRequest(req).(*request.FolderCreateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)

	// create the actual folder
	fold, err := controller.folderService.Create(user, validReq.Name, validReq.OrganizationId, validReq.ParentFolderId)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	// wrap the folder with acl information
	wrapped, err := controller.aclService.Wrap(user, []folder.Folder{*fold})
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("created folder but failed to find user access"))
		return
	}

	// return the data
	util.WriteJsonToResponse(w, http.StatusCreated, wrapped[0])
}

func (controller *FolderController) update(w http.ResponseWriter, req *http.Request) {
	validReq := controller.validatorService.GetModelFromRequest(req).(*request.FolderCreateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	id := chi.URLParam(req, "id")

	fold, err := controller.folderService.Update(user, id, validReq.Name)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	// wrap the folder with acl information
	wrapped, err := controller.aclService.Wrap(user, []folder.Folder{*fold})
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("updated folder but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped)
}

func (controller *FolderController) list(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	organizationId := chi.URLParam(req, "organizationId")
	queryParentFolderId := req.URL.Query().Get("parentFolderId")
	pagination := shared.NewPagination(req)

	var parentFolderId *string
	if len(queryParentFolderId) > 0 {
		parentFolderId = &queryParentFolderId
	}

	folders, err := controller.folderService.List(user, organizationId, parentFolderId, pagination)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	if len(folders) == 0 {
		util.WriteJsonToResponse(w, http.StatusOK, folders)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, folders)
	if err != nil {
		util.WriteHttpError(w, shared.NewInternalServerError("failed to load user access for folders"))
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped)
}