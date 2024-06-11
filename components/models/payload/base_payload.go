package payload

import "time"

type BasePayload struct {
	Search        string
	Limit         *int
	Offset        *int
	LastCreatedAt *time.Time
	Orders        []Orders
}

type Orders struct {
	Order   string
	OrderBy string
}
