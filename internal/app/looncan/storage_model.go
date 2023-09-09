package looncan

import (
	"database/sql"
	"time"
)

type model struct {
	ID         uint64
	Name       string
	Value      string
	ParentID   uint64
	ParentType string
	CreatedAt  time.Time
	UpdatedAt  sql.NullTime
}
