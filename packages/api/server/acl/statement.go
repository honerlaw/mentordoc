package acl

import "github.com/honerlaw/mentordoc/server"

type Statement struct {
	server.Entity

	Effect       string `json:"effect"`
	ResourceName string `json:"resourceName"`
	ResourceID   string `json:"resourceId"`
}
