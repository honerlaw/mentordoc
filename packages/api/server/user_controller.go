package server

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/model"
	"github.com/honerlaw/mentordoc/server/util"
	"net/http"
)

type UserController struct {
	userService           *UserService
	validatorService      *util.ValidatorService
	authenticationService *AuthenticationService
}

func NewUserController(userService *UserService,
	validatorService *util.ValidatorService,
	authenticationService *AuthenticationService) *UserController {
	return &UserController{
		userService:           userService,
		validatorService:      validatorService,
		authenticationService: authenticationService,
	}
}

func (controller *UserController) RegisterRoutes(router chi.Router) {
	router.
		With(controller.validatorService.Middleware(model.UserSigninRequest{})).
		Post("/user/auth", controller.signin)

	router.
		With(controller.validatorService.Middleware(model.UserSignupRequest{})).
		Post("/user", controller.signup)
}

func (controller *UserController) signin(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*model.UserSigninRequest)

	user, err := controller.userService.Authenticate(request.Email, request.Password)
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	accessToken, err := controller.authenticationService.GenerateToken(user.Id, TokenAccess)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.authenticationService.GenerateToken(user.Id, TokenRefresh)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	util.WriteJsonToResponse(w, http.StatusOK, &model.AuthenticationResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}

func (controller *UserController) signup(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*model.UserSignupRequest)

	user, err := controller.userService.Create(request.Email, request.Password);
	if err != nil {
		util.WriteHttpError(w, err)
		return
	}

	accessToken, err := controller.authenticationService.GenerateToken(user.Id, TokenAccess)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.authenticationService.GenerateToken(user.Id, TokenRefresh)
	if err != nil {
		util.WriteHttpError(w, err)
		return;
	}

	util.WriteJsonToResponse(w, http.StatusOK, &model.AuthenticationResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}
