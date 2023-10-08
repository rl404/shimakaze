package cron

import (
	"context"

	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/service"
	"github.com/rl404/shimakaze/internal/utils"
	"github.com/rl404/shimakaze/pkg/log"
)

// Cron contains functions for cron.
type Cron struct {
	service service.Service
}

// New to create new cron.
func New(service service.Service) *Cron {
	return &Cron{
		service: service,
	}
}

func (c *Cron) log(ctx context.Context) {
	errStack := errors.Get(ctx)
	if len(errStack) > 0 {
		utils.Log(map[string]interface{}{
			"level": log.ErrorLevel,
			"error": errStack,
		})
	}
}
