package server

import (
	"github.com/go-chi/chi"
	"net/http"
)

type UserController struct {
	userService *UserService
	validatorService *ValidatorService
}

func NewUserController(userService *UserService, validatorService *ValidatorService) *UserController {
	return &UserController{
		userService: userService,
		validatorService: validatorService,
	}
}

func (controller *UserController) RegisterRoutes(router chi.Router) {
	router.Route("/user/auth", func (r chi.Router) {
		r.Use(controller.validatorService.Middleware(UserSigninRequest{}))
		r.Post("/", controller.signin)
	})
	router.Route("/user", func (r chi.Router) {
		r.Use(controller.validatorService.Middleware(UserSignupRequest{}))
		r.Post("/", controller.signup)
	})
}

func (controller *UserController) signin(w http.ResponseWriter, req *http.Request) {
	model := controller.validatorService.GetModelFromRequest(req).(*UserSigninRequest)

	user, err := controller.userService.Authenticate(model.Email, model.Password)
	if err != nil {
		WriteHttpError(w, err)
		return
	}

	WriteJsonToResponse(w, http.StatusOK, user)
}

func (controller *UserController) signup(w http.ResponseWriter, req *http.Request) {
	model := controller.validatorService.GetModelFromRequest(req).(*UserSignupRequest)

	user, err := controller.userService.Create(model.Email, model.Password);
	if err != nil {
		WriteHttpError(w, err)
		return
	}

	WriteJsonToResponse(w, http.StatusOK, user)
}
