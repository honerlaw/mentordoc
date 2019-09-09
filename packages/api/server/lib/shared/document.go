package shared

type Document struct {
	Entity

	OrganizationId string          `json:"organizationId"`
	FolderId       *string         `json:"folderId"`
	Drafts         []DocumentDraft `json:"drafts"`
}
