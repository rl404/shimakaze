package utils

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/shimakaze/internal/errors"
)

// PubsubRecoverer is custom pubsub recoverer middleware.
func PubsubRecoverer(next pubsub.HandlerFunc) pubsub.HandlerFunc {
	return func(ctx context.Context, message []byte) error {
		defer func() {
			if rvr := recover(); rvr != nil {
				stack.Wrap(
					ctx,
					fmt.Errorf("%s", debug.Stack()),
					fmt.Errorf("%v", rvr),
					errors.ErrInternalServer)
			}
		}()
		return next(ctx, message)
	}
}
