package acl

type AclWrappedModel struct {
	Model   interface{} `json:"model"`
	Actions []string    `json:"actions"`
}
