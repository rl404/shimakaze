package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors"
	"github.com/rl404/shimakaze/internal/utils"
)

// Update to old data.
func (c *Cron) Update(nrApp *newrelic.Application, limit int) error {
	ctx := errors.Init(context.Background())
	defer c.log(ctx)

	tx := nrApp.StartTransaction("Cron update")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueOldAgency(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldActiveVtuber(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	if err := c.queueOldRetiredVtuber(ctx, nrApp, limit); err != nil {
		tx.NoticeError(err)
		return errors.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueOldAgency(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldAgency").End()

	cnt, _, err := c.service.QueueOldAgency(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d agency", cnt)
	nrApp.RecordCustomEvent("QueueOldAgency", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldActiveVtuber(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldActiveVtuber").End()

	cnt, _, err := c.service.QueueOldActiveVtuber(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d active vtuber", cnt)
	nrApp.RecordCustomEvent("QueueOldActiveVtuber", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldRetiredVtuber(ctx context.Context, nrApp *newrelic.Application, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldRetiredVtuber").End()

	cnt, _, err := c.service.QueueOldRetiredVtuber(ctx, limit)
	if err != nil {
		return errors.Wrap(ctx, err)
	}

	utils.Info("queued %d retired vtuber", cnt)
	nrApp.RecordCustomEvent("QueueOldVRetiredtuber", map[string]interface{}{"count": cnt})

	return nil
}
