package request

type SignupRequest struct {
	Email string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	VerifyPassword string `validate:"required,min=8"`
}
