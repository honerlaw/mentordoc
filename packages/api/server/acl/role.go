package acl

import "github.com/honerlaw/mentordoc/server"

type Role struct {
	server.Entity

	Name string `json:"name"`
}