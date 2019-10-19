package request

type DocumentUpdateRequest struct {
	DocumentId    string  `json:"documentId" validate:"required"`
	DraftId       string  `json:"draftId" validate:"required"`
	Name          *string `json:"name"`
	Content       *string `json:"content"`
	ShouldPublish bool    `json:"shouldPublish"`
	ShouldRetract bool    `json:"shouldRetract"`
}
