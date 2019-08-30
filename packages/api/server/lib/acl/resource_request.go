package acl

type ResourceRequest struct {
	ResourcePath []string
	ResourceIds  []string
	Action       *string
}
