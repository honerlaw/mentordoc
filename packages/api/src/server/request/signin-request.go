package request

type SigninRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"reqiured"`
}
