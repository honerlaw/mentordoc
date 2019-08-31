package controller

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/http/middleware"
	"github.com/honerlaw/mentordoc/server/lib/acl"
	"github.com/honerlaw/mentordoc/server/lib/organization"
	"github.com/honerlaw/mentordoc/server/lib/shared"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"log"
	"net/http"
)

type OrganizationController struct {
	organizationService      *organization.OrganizationService
	authenticationMiddleware *middleware.AuthenticationMiddleware
	aclService               *acl.AclService
}

func NewOrganizationController(
	organizationService *organization.OrganizationService,
	authenticationMiddleware *middleware.AuthenticationMiddleware,
	aclService *acl.AclService,
) *OrganizationController {
	return &OrganizationController{
		organizationService:      organizationService,
		authenticationMiddleware: authenticationMiddleware,
		aclService:               aclService,
	}
}

func (controller *OrganizationController) RegisterRoutes(router chi.Router) {
	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/organization/list", controller.list)
}

func (controller *OrganizationController) list(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)

	orgs, err := controller.organizationService.List(user)
	if err != nil {
		log.Print(err)
		util.WriteHttpError(w, err)
		return
	}

	if len(orgs) == 0 {
		log.Print("no orgs")
		util.WriteJsonToResponse(w, http.StatusOK, orgs)
		return
	}

	wrapped, err := controller.aclService.Wrap(user, orgs)
	if err != nil {
		log.Print(err)
		util.WriteHttpError(w, shared.NewInternalServerError("found organizations but failed to find user access"))
		return
	}

	util.WriteJsonToResponse(w, http.StatusOK, wrapped)
}
