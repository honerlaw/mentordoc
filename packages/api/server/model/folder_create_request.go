package model

type FolderCreateRequest struct {
	OrganizationId string  `json:"organizationId" validate:"required"`
	Name           string  `json:"name" validate:"required"`
	ParentFolderId *string `json:"parentFolderId"`
}
