package server

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	"net/http"
)

type FolderController struct {
	validatorService         *util.ValidatorService
	folderService            *FolderService
	authenticationMiddleware *AuthenticationMiddleware
}

func NewFolderController(validatorService *util.ValidatorService, folderService *FolderService, authenticationMiddleware *AuthenticationMiddleware) *FolderController {
	return &FolderController{
		validatorService:         validatorService,
		folderService:            folderService,
		authenticationMiddleware: authenticationMiddleware,
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

	folder, err := controller.folderService.Create(user, request.Name, request.OrganizationId, request.ParentFolderId)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	// todo we need to merge acl actions with the folder
}

func (controller *FolderController) list(w http.ResponseWriter, req *http.Request) {
	// do nothing
}
