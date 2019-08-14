package server

type User struct {
	Entity

	Email string `json:"email"`

	Password string
}
