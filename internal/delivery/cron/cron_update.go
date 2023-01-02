package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/shimakaze/internal/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// Update to old data.
func (c *Cron) Update(nrApp *newrelic.Application, limit int) error {
	ctx := errors.Init(context.Background())
	defer c.log(ctx)

	tx := nrApp.StartTransaction("Cron update")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueOldAgency(ctx, nrApp); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldVtuber(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueOldAgency(ctx context.Context, nrApp *newrelic.Application) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldAgency").End()

	cnt, _, err := c.service.QueueOldAgency(ctx)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d agency", cnt)
	nrApp.RecordCustomEvent("QueueOldAgency", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldVtuber(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldVtuber").End()

	cnt, _, err := c.service.QueueOldVtuber(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d vtuber", cnt)
	nrApp.RecordCustomEvent("QueueOldVtuber", map[string]interface{}{"count": cnt})

	return nil
}
