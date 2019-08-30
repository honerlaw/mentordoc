package document

import "github.com/honerlaw/mentordoc/server/lib/shared"

type DocumentContent struct {
	shared.Entity

	DocumentId string `json:"documentId"`
	Content    string `json:"content"`
}
