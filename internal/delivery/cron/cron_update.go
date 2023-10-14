package cron

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/shimakaze/internal/utils"
)

// Update to old data.
func (c *Cron) Update(limit int) error {
	ctx := stack.Init(context.Background())
	defer c.log(ctx)

	tx := c.nrApp.StartTransaction("Cron update")
	defer tx.End()

	ctx = newrelic.NewContext(ctx, tx)

	if err := c.queueOldAgency(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	if err := c.queueOldActiveVtuber(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	if err := c.queueOldRetiredVtuber(ctx, limit); err != nil {
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (c *Cron) queueOldAgency(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldAgency").End()

	cnt, _, err := c.service.QueueOldAgency(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d agency", cnt)
	c.nrApp.RecordCustomEvent("QueueOldAgency", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldActiveVtuber(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldActiveVtuber").End()

	cnt, _, err := c.service.QueueOldActiveVtuber(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d active vtuber", cnt)
	c.nrApp.RecordCustomEvent("QueueOldActiveVtuber", map[string]interface{}{"count": cnt})

	return nil
}

func (c *Cron) queueOldRetiredVtuber(ctx context.Context, limit int) error {
	defer newrelic.FromContext(ctx).StartSegment("queueOldRetiredVtuber").End()

	cnt, _, err := c.service.QueueOldRetiredVtuber(ctx, limit)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	utils.Info("queued %d retired vtuber", cnt)
	c.nrApp.RecordCustomEvent("QueueOldVRetiredtuber", map[string]interface{}{"count": cnt})

	return nil
}
