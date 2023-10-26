package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
)

// Cron contains functions for cron.
type Cron struct {
	service service.Service
	nrApp   *newrelic.Application
}

// New to create new cron.
func New(service service.Service, nrApp *newrelic.Application) *Cron {
	return &Cron{
		service: service,
		nrApp:   nrApp,
	}
}

func (c *Cron) log(ctx context.Context) {
	errStack := stack.Get(ctx)
	if len(errStack) > 0 {
		utils.Log(map[string]interface{}{
			"level": utils.ErrorLevel,
			"error": errStack,
		})
	}
}
