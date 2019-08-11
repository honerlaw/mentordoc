package controller

import (
	"github.com/go-chi/chi"
	"net/http"
	"server/service"
)

type User struct {
	userService *service.User
}

func NewUser(userService *service.User) *User {
	return &User{
		userService: userService,
	}
}

func (controller *User) RegisterRoutes(router chi.Router) {
	router.Post("/signin", controller.signin)
	router.Post("/signup", controller.signup)
}

func (controller *User) signin(w http.ResponseWriter, r *http.Request) {

}

func (controller *User) signup(w http.ResponseWriter, r *http.Request) {

}
