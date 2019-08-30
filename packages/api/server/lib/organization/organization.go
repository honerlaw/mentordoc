package organization

import "github.com/honerlaw/mentordoc/server/lib/shared"

type Organization struct {
	shared.Entity

	Name string `json:"name"`
}
