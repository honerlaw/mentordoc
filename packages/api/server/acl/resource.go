package acl

type ResourceRequest struct {
	ResourcePath []string
	ResourceIds  []string
	Action       *string
}

type ResourceResponse struct {
	PermissionId string
	UserId       string
	ResourcePath string
	ResourceId   string
	Action       string
}
