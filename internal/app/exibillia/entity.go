package exibillia

import (
	"time"
)

// Exibillia is minimalistic entity without any business logic, only simple CRUD
type Exibillia struct {
	ID          uint64
	Name        string
	Description string
	Tags        []string

	CreatedAt time.Time
	UpdatedAt time.Time
}
