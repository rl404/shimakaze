package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/utils"
)

// Fill to fill missing data.
func (c *Cron) Fill(nrApp *newrelic.Application, limit int) error {
	ctx := stack.Init(context.Background())
	defer c.log(ctx)

	tx := nrApp.StartTransaction("Cron fill")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueMissingAgency(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return stack.Wrap(ctx, err)
	}

	if err := c.queueMissingVtuber(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueMissingAgency(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueMissingAgency").End()

	cnt, _, err := c.service.QueueMissingAgency(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d agency", cnt)
	nrApp.RecordCustomEvent("QueueMissingAgency", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueMissingVtuber(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueMissingVtuber").End()

	cnt, _, err := c.service.QueueMissingVtuber(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d vtuber", cnt)
	nrApp.RecordCustomEvent("QueueMissingVtuber", map[string]interface{}{"count": cnt})

	return nil
}
