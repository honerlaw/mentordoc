package model

type DocumentUpdateRequest struct {
	DocumentId string `json:"documentId" validate:"required''"`
	Name       string `json:"name" validate:"required"`
	Content    string `json:"content" validate:"required"`
}
