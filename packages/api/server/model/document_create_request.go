package model

type DocumentCreateRequest struct {
	OrganizationId string `json:"organizationId" validate:"required"`
	FolderId       *string `json:"folderId"`
	Name           string `json:"name" validate:"required"`
	Content        string `json:"content" validate:"required"`
}
