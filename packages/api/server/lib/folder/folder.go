package folder

import "github.com/honerlaw/mentordoc/server/lib/shared"

type Folder struct {
	shared.Entity

	Name           string  `json:"name"`
	OrganizationId string  `json:"organizationId"`
	ParentFolderId *string `json:"parentFolderId"`
}
