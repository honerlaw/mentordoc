package controller

import (
	"github.com/go-chi/chi"
	"net/http"
	"server/request"
	"server/service"
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
		r.Use(controller.validatorService.Middleware(&request.SignupRequest{}))
		r.Post("/", controller.signin)
	})
	router.Route("/signup", func (r chi.Router) {
		r.Use(controller.validatorService.Middleware(&request.SignupRequest{}))
		r.Post("/", controller.signin)
	})
}

func (controller *User) signin(w http.ResponseWriter, req *http.Request) {
	controller.validatorService.GetModelFromRequest(req)
}

func (controller *User) signup(w http.ResponseWriter, req *http.Request) {
	controller.validatorService.GetModelFromRequest(req)
}
