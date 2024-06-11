package request

import (
	"time"
)

type Order struct {
	Order   string `json:"order"`
	OrderBy string `json:"orderBy"`
}

type BaseRequest struct {
	Id        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type BaseGetListParams struct {
	Search string `json:"search"`
	// Orders []Order `json:"orders"`
	Limit  int `json:"limit" validate:"gte=0"`
	Offset int `json:"offset" validate:"gte=0"`
}
