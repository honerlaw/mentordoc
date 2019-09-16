package request

type DocumentUpdateRequest struct {
	DocumentId string `json:"documentId" validate:"required"`
	DraftId    string `json:"draftId" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Content    string `json:"content" validate:"required"`
}
