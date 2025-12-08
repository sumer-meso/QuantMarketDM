package rabbitmq

import (
	"context"
	"time"

	"github.com/sumer-meso/QuantMarketDM/utils/logging"
)

func timeSleepSelect(ctx context.Context, ms int) {
	select {
	case <-ctx.Done():
		logging.Logf("[RabbitMq][Base] ctx cancelled outside. timeSleepSelect closing.")
		return
	case <-time.After(time.Duration(ms) * time.Millisecond):
	}
}
