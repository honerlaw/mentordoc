package shared

type DocumentContent struct {
	Entity

	DocumentId string `json:"documentId"`
	Content    string `json:"content"`
}