package wiredata

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type WSEventBase struct {
	Event string `json:"e"`
	Time  int64  `json:"E"`
}

type EventBase struct {
	Event string
	Time  int64
}

type LocalBase struct {
	LocalTime int64
}

type MessageOverRabbitMQ struct {
	RoutingKey     string
	DataIdentifier string
	Body           []byte
}

func (m *MessageOverRabbitMQ) PublishOnWire(ctx context.Context, ch *amqp.Channel, exchange string, expiration string) error {
	return ch.PublishWithContext(ctx,
		exchange, m.RoutingKey, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Headers:      amqp.Table{"x-msg-type": m.DataIdentifier},
			Expiration:   expiration,
			DeliveryMode: amqp.Transient,
			Body:         m.Body,
		})
}

func (m *MessageOverRabbitMQ) RetrieveFromWire(ctx context.Context, del *amqp.Delivery) error {
	if v, ok := del.Headers["x-msg-type"]; ok {
		if s, ok := v.(string); ok {
			m.DataIdentifier = s
		}
	}
	m.Body = del.Body
	m.RoutingKey = del.RoutingKey
	return nil
}

type NotMatchError struct {
	Expected string // e.g., "orderbook"
	Actual   string // e.g., "trade"
}

func (e NotMatchError) Error() string {
	return fmt.Sprintf("data identifier not match, expected: %s, actual: %s", e.Expected, e.Actual)
}

type RMQRoutingIdentifier interface{ RMQRoutingIdentifier() string }
type RMQDataIdentifier interface{ RMQDataIdentifier() string }

type RMQSerilizationOnWire interface {
	RMQEncodeMessage() (MessageOverRabbitMQ, error)
	RMQDecodeMessage(m MessageOverRabbitMQ) error
}

var _ = []interface {
	RMQRoutingIdentifier
	RMQDataIdentifier
	RMQSerilizationOnWire
}{
	(*Trade)(nil),
	(*Kline)(nil),
	(*OrderBook)(nil),
}
