package model

type DocumentContent struct {
	Entity

	DocumentId string `json:"documentId"`
	Content    string `json:"content"`
}
