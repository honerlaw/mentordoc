package model

type Document struct {
	Entity

	Name           string           `json:"name"`
	OrganizationId string           `json:"organizationId"`
	FolderId       *string          `json:"folderId"`
	Content        *DocumentContent `json:"content,omitempty"`
}
