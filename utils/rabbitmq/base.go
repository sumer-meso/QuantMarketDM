package rabbitmq

import (
	"context"
	"time"
)

func timeSleepSelect(ctx context.Context, ms int) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Duration(ms) * time.Millisecond):
	}
}
