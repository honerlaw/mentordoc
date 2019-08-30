package acl

import "github.com/honerlaw/mentordoc/server/lib/shared"

type Permission struct {
	shared.Entity

	ResourcePath string `json:"resourcePath"`
	Action       string `json:"action"`
}
