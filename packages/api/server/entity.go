package server

type Entity struct {
	Id string `json:"id"`

	UpdatedAt int64 `json:"updatedAt"`

	CreatedAt int64 `json:"createdAt"`

	DeletedAt *int64 `json:"deletedAt"` // nil if actualy deleted
}
