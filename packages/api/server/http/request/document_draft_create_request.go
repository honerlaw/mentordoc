package request

type DocumentDraftCreateRequest struct {
	DocumentId string `json:"documentId" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Content    string `json:"content" validate:"required"`
}
