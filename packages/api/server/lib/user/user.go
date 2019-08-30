package user

import (
	"github.com/honerlaw/mentordoc/server/lib/shared"
)

type User struct {
	shared.Entity

	Email string `json:"email"`

	Password string `json:"-"`
}
