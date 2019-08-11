package controller

import (
	"github.com/go-chi/chi"
	"net/http"
	"server/request"
	"server/service"
	"server/util"
)

type User struct {
	userService *service.User
	validatorService *service.Validator
}

func NewUser(userService *service.User, validatorService *service.Validator) *User {
	return &User{
		userService: userService,
		validatorService: validatorService,
	}
}

func (controller *User) RegisterRoutes(router chi.Router) {
	router.Route("/signin", func (r chi.Router) {
		r.Use(controller.validatorService.Middleware(request.SigninRequest{}))
		r.Post("/", controller.signin)
	})
	router.Route("/signup", func (r chi.Router) {
		r.Use(controller.validatorService.Middleware(request.SignupRequest{}))
		r.Post("/", controller.signin)
	})
}

func (controller *User) signin(w http.ResponseWriter, req *http.Request) {
	model := controller.validatorService.GetModelFromRequest(req).(*request.SigninRequest)

	user, err := controller.userService.Authenticate(model.Email, model.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	util.WriteJsonToResponse(w, user)
}

func (controller *User) signup(w http.ResponseWriter, req *http.Request) {
	model := controller.validatorService.GetModelFromRequest(req).(*request.SignupRequest)

	user, err := controller.userService.Create(model.Email, model.Password);
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	util.WriteJsonToResponse(w, user)
}
