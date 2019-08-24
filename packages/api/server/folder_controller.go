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
	router.Get("/folder", controller.list)
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
	wrapped, err := controller.aclService.Wrap(user, []*model.Folder{folder})
	if err != nil {
		util.WriteHttpError(w, model.NewInternalServerError("created folder but failed to find user access"))
		return
	}

	// return the data
	util.WriteJsonToResponse(w, http.StatusCreated, wrapped[0])
}

func (controller *FolderController) list(w http.ResponseWriter, req *http.Request) {
	// do nothing
}
