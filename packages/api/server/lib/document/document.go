package document

import "github.com/honerlaw/mentordoc/server/lib/shared"

type Document struct {
	shared.Entity

	Name           string           `json:"name"`
	OrganizationId string           `json:"organizationId"`
	FolderId       *string          `json:"folderId"`
	Content        *DocumentContent `json:"content,omitempty"`
}
