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
}

func (controller *UserController) get(w http.ResponseWriter, req *http.Request) {
	user := controller.authenticationMiddleware.GetUserFromRequest(req)

	util.WriteJsonToResponse(w, http.StatusOK, user)
}

func (controller *UserController) signin(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*request.UserSigninRequest)

	user, err := controller.userService.Authenticate(request.Email, request.Password)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	accessToken, err := controller.tokenService.GenerateToken(user.Id, util.TokenAccess)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.tokenService.GenerateToken(user.Id, util.TokenRefresh)
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
	request := controller.validatorService.GetModelFromRequest(req).(*request.UserSignupRequest)

	user, err := controller.userService.Create(request.Email, request.Password);
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	accessToken, err := controller.tokenService.GenerateToken(user.Id, util.TokenAccess)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.tokenService.GenerateToken(user.Id, util.TokenRefresh)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	util.WriteJsonToResponse(w, http.StatusOK, &response.AuthenticationResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}
