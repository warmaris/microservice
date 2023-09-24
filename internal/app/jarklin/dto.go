package jarklin

import (
	"time"
)

type NotifyStatus uint8

const (
	StatusUnknown NotifyStatus = iota
	StatusPrepared
	StatusSent
	StatusSuccess
	StatusFail
)

type NotifyRequest struct {
	MessageID  string    `json:"message_id"`
	EntityID   uint64    `json:"entity_id"`
	EntityName string    `json:"entity_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type NotifyResponse struct {
	MessageID  string       `json:"message_id"`
	Status     NotifyStatus `json:"status"`
	StatusInfo string       `json:"status_info"`
}
