package cron

import (
	"context"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type JobName string
type JobSchedule string
type JobConfig map[JobName]JobSchedule

const (
	JarklinSendNewEvents JobName = "jarklin.send_new_events"
)

type JarklinService interface {
	SendNewEvents(ctx context.Context)
}

type Cron struct {
	runner *cron.Cron
	jarklinService JarklinService
}

func NewCron(config JobConfig, jarklin JarklinService) *Cron {
	c := new(Cron)
	runner := cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger)))
	for name, schedule := range config {
		job := c.getMethodByName(name)
		if job == nil {
			zap.S().Warnw("no method for job", "job", name)
			
			continue
		}
		runner.AddFunc(string(schedule), job)
	}

	c.runner = runner
	c.jarklinService = jarklin

	return c
}

func (c *Cron) Start() {
	c.runner.Start()
}

func (c *Cron) Stop() {
	ctx := c.runner.Stop()
	<- ctx.Done()
}

func (c *Cron) JarklinSendNewEvents() {
	c.jarklinService.SendNewEvents(context.Background())
}

func (c *Cron) getMethodByName(name JobName) func() {
	if name == JarklinSendNewEvents {
		return c.JarklinSendNewEvents
	}

	return nil
}