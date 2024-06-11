package payload

import (
	"time"

	"github.com/google/uuid"
)

type GetMenusPayload struct {
	BasePayload
	Ids []uuid.UUID
}
type DeleteMenusPayload struct {
	Ids       []uuid.UUID
	DeletedAt *time.Time
}
