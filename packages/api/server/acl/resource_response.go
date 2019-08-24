package acl

type ResourceResponse struct {
	PermissionId string
	UserId       string
	ResourcePath string
	ResourceId   string
	Action       string
}

type ResourcePaginatedResponse struct {
	ResourceId string
	Actions []string
}
