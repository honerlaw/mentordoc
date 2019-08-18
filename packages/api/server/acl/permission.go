package acl

import "github.com/honerlaw/mentordoc/server"

type Permission struct {
	server.Entity

	ResourcePath string `json:"resourcePath"`
	Action       string `json:"action"`
}
