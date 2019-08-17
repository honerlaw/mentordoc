package acl

import "github.com/honerlaw/mentordoc/server"

type Statement struct {
	server.Entity

	ResourceName string `json:"resourceName"`
	ResourceID   string `json:"resourceId"`
	Action       string `json:"action"`
}
