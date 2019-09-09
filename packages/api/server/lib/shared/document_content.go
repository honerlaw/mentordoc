package shared

type DocumentContent struct {
	Entity

	DocumentDraftId string `json:"documentDraftId"`
	Content    string `json:"content"`
}
