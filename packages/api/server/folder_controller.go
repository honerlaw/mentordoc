package server

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/acl"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	"net/http"
)

type FolderController struct {
	validatorService         *util.ValidatorService
	folderService            *FolderService
	authenticationMiddleware *AuthenticationMiddleware
	aclService               *acl.AclService
}

func NewFolderController(
	validatorService *util.ValidatorService,
	folderService *FolderService,
	authenticationMiddleware *AuthenticationMiddleware,
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
		With(controller.validatorService.Middleware(model.FolderCreateRequest{}),
			controller.authenticationMiddleware.HasAccessToken()).
		Post("/folder", controller.create)
	router.
		With(controller.validatorService.Middleware(model.FolderUpdateRequest{}),
			controller.authenticationMiddleware.HasAccessToken()).
		Put("/folder/{id}", controller.update)
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/folder/list/{organizationId}", controller.list)
}

func (controller *FolderController) create(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*model.FolderCreateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)

	// create the actual folder
	folder, err := controller.folderService.Create(user, request.Name, request.OrganizationId, request.ParentFolderId)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	// wrap the folder with acl information
	wrapped, err := controller.aclService.Wrap(user, []model.Folder{*folder})
	if err != nil {
		util.WriteHttpError(w, model.NewInternalServerError("created folder but failed to find user access"))
		return
	}

	// return the data
	util.WriteJsonToResponse(w, http.StatusCreated, wrapped[0])
}

func (controller *FolderController) update(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*model.FolderCreateRequest)
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	id := chi.URLParam(req, "id")

	folder, err := controller.folderService.Update(user, id, request.Name)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	// wrap the folder with acl information
	wrapped, err := controller.aclService.Wrap(user, []model.Folder{*folder})
	if err != nil {
		util.WriteHttpError(w, model.NewInternalServerError("updated folder but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped)
}

func (controller *FolderController) list(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)
	organizationId := chi.URLParam(req, "organizationId")
	queryParentFolderId := req.URL.Query().Get("parentFolderId")
	pagination := model.NewPagination(req)

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
		util.WriteHttpError(w, model.NewInternalServerError("failed to load user access for folders"))
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped)
}