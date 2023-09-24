// Package jarklin provides creation of entity with event dispatching.
// There are three ways of dispatching event:
// - Sync (while request processing)
// - Async (using background job)
// - Combined (try to send within request, retry in background later if error occurs)
// Also it shows how to integrate MQ in project and use it for async Send/Reply communication.
package jarklin

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

type Storage interface {
	save(ctx context.Context, entity *Jarklin) error
	updateStatus(ctx context.Context, msgID string, status NotifyStatus, statusInfo string) error
	getNewEvents(ctx context.Context) ([]NotifyRequest, error)
}

type Sender interface {
	send(ctx context.Context, ev NotifyRequest) error
}

type Service struct {
	storage Storage
	sender  Sender
}

func NewService(storage Storage, sender Sender) *Service {
	return &Service{storage: storage, sender: sender}
}

// CreateAndSend dispatches event syncronously.
// Prepared event is saved with entity in transaction, so we can retry send if error occurs.
// If we cannot save entity in db, nothing will be sent.
func (s *Service) CreateAndSend(ctx context.Context, name string) error {
	entity := &Jarklin{
		Name:      name,
		CreatedAt: time.Now(),
	}

	err := s.storage.save(ctx, entity)
	if err != nil {
		zap.S().Errorw("storage create", "error", err)

		return err
	}

	err = s.sender.send(ctx, entity.GetEvent())
	if err != nil {
		zap.S().Errorw("send event", "error", err)

		return err
	}

	err = s.storage.updateStatus(ctx, entity.GetEvent().MessageID, StatusSent, "")
	if err != nil {
		zap.S().Errorw("update status", "error", err)

		return err
	}

	return nil
}

// CreateAndSave just saves prepared message with entity in transaction. 
// Events will be dispatched later from cron job.
func (s *Service) CreateAndSave(ctx context.Context, name string) error {
	entity := &Jarklin{
		Name:      name,
		CreatedAt: time.Now(),
	}

	err := s.storage.save(ctx, entity)
	if err != nil {
		zap.S().Errorw("storage create", "error", err)
	}

	return err
}

// SendNewEvents is called from cron job and sends prepared and not sent events.
func (s *Service) SendNewEvents(ctx context.Context) {
	events, err := s.storage.getNewEvents(ctx)
	if err != nil {
		zap.S().Errorw("get new events from storage", "error", err)

		return
	}

	for _, ev := range events {
		err = s.sender.send(ctx, ev)
		if err != nil {
			zap.S().Errorw("send event", "error", err)

			return
		}

		s.storage.updateStatus(ctx, ev.MessageID, StatusSent, "")
	}

	zap.S().Infow("new events sent", "count", len(events))
}

// HandleResponse updates status for dispatched events, processed by another microservice.
func (s *Service) HandleResponse(ctx context.Context, msg []byte) error {
	res := new(NotifyResponse)
	err := json.Unmarshal(msg, res)
	if err != nil {
		zap.S().Errorw("unmarshal json", "error", err)

		return err
	}

	err = s.storage.updateStatus(ctx, res.MessageID, res.Status, res.StatusInfo)
	if err != nil {
		zap.S().Errorw("update status from response", "error", err)

		return err
	}

	return nil
}
