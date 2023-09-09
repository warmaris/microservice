package looncan

import "time"

type ParentType string

const (
	ParentTypeAcaer     ParentType = "acaer"
	ParentTypeExibillia ParentType = "exibillia"
)

type Looncan struct {
	ID         uint64
	Name       string
	Value      string
	ParentID   uint64
	ParentType ParentType
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
