package wiredata

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const HeaderKeyType = "x-tide-type"
const HeaderKeyTable = "x-tide-table"
const HeaderKeyIndex = "x-tide-index"

const defaultMsgTTLInQueue = "20000"

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
	Source    string
}

type MessageOverRabbitMQ struct {
	RoutingKey     string
	DataIdentifier string
	StoreTable     string
	StoreIndex     string
	Body           []byte
}

func (m *MessageOverRabbitMQ) PublishOnWire(ctx context.Context, ch *amqp.Channel, exchange string) error {
	return ch.PublishWithContext(ctx,
		exchange, m.RoutingKey, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Headers: amqp.Table{
				HeaderKeyType:  m.DataIdentifier,
				HeaderKeyTable: m.StoreTable,
				HeaderKeyIndex: m.StoreIndex,
			},
			Expiration:   defaultMsgTTLInQueue,
			DeliveryMode: amqp.Transient,
			Body:         m.Body,
		})
}

func (m *MessageOverRabbitMQ) RetrieveFromWire(ctx context.Context, del *amqp.Delivery) error {
	if v, ok := del.Headers[HeaderKeyType]; ok {
		if s, ok := v.(string); ok {
			m.DataIdentifier = s
		}
	}
	if v, ok := del.Headers[HeaderKeyTable]; ok {
		if s, ok := v.(string); ok {
			m.StoreTable = s
		}
	}
	if v, ok := del.Headers[HeaderKeyIndex]; ok {
		if s, ok := v.(string); ok {
			m.StoreIndex = s
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
type RMQDataStoreTable interface{ RMQDataStoreTable() string }
type RMQDataStoreIndex interface{ RMQDataStoreIndex() string }

type RMQSerilizationOnWire interface {
	RMQEncodeMessage() (MessageOverRabbitMQ, error)
	RMQDecodeMessage(m MessageOverRabbitMQ) error
}

var _ = []interface {
	RMQRoutingIdentifier
	RMQDataIdentifier
	RMQDataStoreTable
	RMQDataStoreIndex
	RMQSerilizationOnWire
}{
	(*Trade)(nil),
	(*Kline)(nil),
	(*OrderBook)(nil),
	(*TradeLite)(nil),
	(*AccountUpdate)(nil),
	(*OrderUpdate)(nil),
}
