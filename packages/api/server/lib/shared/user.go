package shared

type User struct {
	Entity

	Email string `json:"email"`

	Password string `json:"-"`
}
