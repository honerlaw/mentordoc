package acl

import "github.com/honerlaw/mentordoc/server/lib/shared"

type Role struct {
	shared.Entity

	Name string `json:"name"`
}