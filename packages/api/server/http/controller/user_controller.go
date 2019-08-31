package controller

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/http/middleware"
	"github.com/honerlaw/mentordoc/server/http/request"
	"github.com/honerlaw/mentordoc/server/http/response"
	"github.com/honerlaw/mentordoc/server/lib/user"
	"github.com/honerlaw/mentordoc/server/lib/util"
	"net/http"
)

type UserController struct {
	userService              *user.UserService
	validatorService         *util.ValidatorService
	tokenService    *util.TokenService
	authenticationMiddleware *middleware.AuthenticationMiddleware
}

func NewUserController(
	userService *user.UserService,
	validatorService *util.ValidatorService,
	tokenService *util.TokenService,
	authenticationMiddleware *middleware.AuthenticationMiddleware,
) *UserController {
	return &UserController{
		userService:              userService,
		validatorService:         validatorService,
		tokenService:    tokenService,
		authenticationMiddleware: authenticationMiddleware,
	}
}

func (controller *UserController) RegisterRoutes(router chi.Router) {
	router.
		With(controller.validatorService.Middleware(request.UserSigninRequest{})).
		Post("/user/auth", controller.signin)

	router.
		With(controller.validatorService.Middleware(request.UserSignupRequest{})).
		Post("/user", controller.signup)

	router.
		With(controller.authenticationMiddleware.HasAccessToken()).
		Get("/user", controller.get)

	router.
		With(controller.authenticationMiddleware.HasRefreshToken()).
		Post("/user/auth/refresh", controller.refreshToken)
}

func (controller *UserController) get(w http.ResponseWriter, req *http.Request) {
	u := controller.authenticationMiddleware.GetUserFromRequest(req)

	util.WriteJsonToResponse(w, http.StatusOK, u)
}

func (controller *UserController) refreshToken(w http.ResponseWriter, req *http.Request) {
	u := controller.authenticationMiddleware.GetUserFromRequest(req)

	accessToken, err := controller.tokenService.GenerateToken(u.Id, util.TokenAccess)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.tokenService.GenerateToken(u.Id, util.TokenRefresh)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	util.WriteJsonToResponse(w, http.StatusOK, &response.AuthenticationResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}

func (controller *UserController) signin(w http.ResponseWriter, req *http.Request) {
	validReq := controller.validatorService.GetModelFromRequest(req).(*request.UserSigninRequest)

	u, err := controller.userService.Authenticate(validReq.Email, validReq.Password)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	accessToken, err := controller.tokenService.GenerateToken(u.Id, util.TokenAccess)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.tokenService.GenerateToken(u.Id, util.TokenRefresh)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	util.WriteJsonToResponse(w, http.StatusOK, &response.AuthenticationResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}

func (controller *UserController) signup(w http.ResponseWriter, req *http.Request) {
	validReq := controller.validatorService.GetModelFromRequest(req).(*request.UserSignupRequest)

	u, err := controller.userService.Create(validReq.Email, validReq.Password);
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	accessToken, err := controller.tokenService.GenerateToken(u.Id, util.TokenAccess)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.tokenService.GenerateToken(u.Id, util.TokenRefresh)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	util.WriteJsonToResponse(w, http.StatusOK, &response.AuthenticationResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}
