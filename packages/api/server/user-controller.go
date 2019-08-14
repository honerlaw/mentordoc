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
	router.Route("/signin", func (r chi.Router) {
		r.Use(controller.validatorService.Middleware(UserSigninRequest{}))
		r.Post("/", controller.signin)
	})
	router.Route("/signup", func (r chi.Router) {
		r.Use(controller.validatorService.Middleware(UserSignupRequest{}))
		r.Post("/", controller.signin)
	})
}

func (controller *UserController) signin(w http.ResponseWriter, req *http.Request) {
	model := controller.validatorService.GetModelFromRequest(req).(*UserSigninRequest)

	user, err := controller.userService.Authenticate(model.Email, model.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	WriteJsonToResponse(w, user)
}

func (controller *UserController) signup(w http.ResponseWriter, req *http.Request) {
	model := controller.validatorService.GetModelFromRequest(req).(*UserSignupRequest)

	user, err := controller.userService.Create(model.Email, model.Password);
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	WriteJsonToResponse(w, user)
}
