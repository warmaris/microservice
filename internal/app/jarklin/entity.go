package jarklin

import (
	"time"

	"github.com/gofrs/uuid"
)

type Jarklin struct {
	ID uint64
	Name string
	CreatedAt time.Time

	event *NotifyRequest
}

func (j *Jarklin) GetEvent() NotifyRequest {
	if j.event == nil {
		newID, _ := uuid.NewV4()
		j.event = &NotifyRequest{
			MessageID: newID.String(),
			EntityID: j.ID,
			EntityName: j.Name,
			CreatedAt: j.CreatedAt,
		}
	}

	return *j.event
}