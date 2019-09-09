package shared

type DocumentDraft struct {
	Entity

	DocumentId  string           `json:"documentId"`
	Name        string           `json:"name"`
	Content     *DocumentContent `json:"content,omitempty"`
	PublishedAt *int64           `json:"publishedAt"`
	RetractedAt *int64           `json:"retractedAt"`
}
