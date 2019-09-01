package shared

type Folder struct {
	Entity

	Name           string  `json:"name"`
	OrganizationId string  `json:"organizationId"`
	ParentFolderId *string `json:"parentFolderId"`
	ChildCount     int     `json:"childCount"`
}
