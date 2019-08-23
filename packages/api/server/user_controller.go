package server

import (
	"github.com/go-chi/chi"
	"github.com/honerlaw/mentordoc/server/model"
	"net/http"
)

type UserController struct {
	userService           *UserService
	validatorService      *ValidatorService
	authenticationService *AuthenticationService
}

func NewUserController(userService *UserService,
	validatorService *ValidatorService,
	authenticationService *AuthenticationService) *UserController {
	return &UserController{
		userService:           userService,
		validatorService:      validatorService,
		authenticationService: authenticationService,
	}
}

func (controller *UserController) RegisterRoutes(router chi.Router) {
	router.Route("/user/auth", func(r chi.Router) {
		r.Use(controller.validatorService.Middleware(UserSigninRequest{}))
		r.Post("/", controller.signin)
	})
	router.Route("/user", func(r chi.Router) {
		r.Use(controller.validatorService.Middleware(UserSignupRequest{}))
		r.Post("/", controller.signup)
	})
}

func (controller *UserController) signin(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*UserSigninRequest)

	user, err := controller.userService.Authenticate(request.Email, request.Password)
	if err != nil {
		WriteHttpError(w, err)
		return
	}

	accessToken, err := controller.authenticationService.GenerateToken(user.Id, TokenAccess)
	if err != nil {
		WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.authenticationService.GenerateToken(user.Id, TokenRefresh)
	if err != nil {
		WriteHttpError(w, err)
		return;
	}

	WriteJsonToResponse(w, http.StatusOK, &model.AuthenticationResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}

func (controller *UserController) signup(w http.ResponseWriter, req *http.Request) {
	request := controller.validatorService.GetModelFromRequest(req).(*UserSignupRequest)

	user, err := controller.userService.Create(request.Email, request.Password);
	if err != nil {
		WriteHttpError(w, err)
		return
	}

	accessToken, err := controller.authenticationService.GenerateToken(user.Id, TokenAccess)
	if err != nil {
		WriteHttpError(w, err)
		return;
	}

	refreshToken, err := controller.authenticationService.GenerateToken(user.Id, TokenRefresh)
	if err != nil {
		WriteHttpError(w, err)
		return;
	}

	WriteJsonToResponse(w, http.StatusOK, &model.AuthenticationResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	})
}
