package shared

type ResourceHistory struct {
	Entity

	ResourceId   string `json:"resourceId"`
	ResourceName string `json:"resourceName"`
	UserId       string `json:"userId"`
	Action       string `json:"action"`
}
