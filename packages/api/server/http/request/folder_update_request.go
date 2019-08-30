package request

type FolderUpdateRequest struct {
	Name string `json:"name" validate:"required"`
}
