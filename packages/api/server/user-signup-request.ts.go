package server

type UserSignupRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	VerifyPassword string `json:"verifyPassword" validate:"required,min=8"`
}
