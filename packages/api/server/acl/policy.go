package acl

import "github.com/honerlaw/mentordoc/server"

type Policy struct {
	server.Entity

	Name string `json:"name"`
}