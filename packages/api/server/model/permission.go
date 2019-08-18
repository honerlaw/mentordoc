package model

type Permission struct {
	Entity

	ResourcePath string `json:"resourcePath"`
	Action       string `json:"action"`
}
