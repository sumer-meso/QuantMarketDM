package proto

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

const HeaderKeyType = "x-tide-type"
const HeaderKeyTable = "x-tide-table"
const HeaderKeyIndex = "x-tide-index"

const UnknownDataIdentifier = "tide.unknown"
const UnknownRoutingIdentifier = "tide.proto.unknown"

const defaultMsgTTLInQueue = "20000"

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
	m.DataIdentifier = UnknownDataIdentifier
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

type Unknown struct {
	RkFromRabbitMQ string
	Body           []byte
}

func (u *Unknown) RMQRoutingIdentifier() string {
	return UnknownRoutingIdentifier
}

func (u *Unknown) RMQDataIdentifier() string {
	return UnknownDataIdentifier
}

func (u *Unknown) RMQEncodeMessage() (*MessageOverRabbitMQ, error) {
	panic("Unknown.RMQEncodeMessage() should never be called")
}

func (u *Unknown) RMQDecodeMessage(m *MessageOverRabbitMQ) error {
	// we are not checking DataIdentifier here, as Unknown is a fallback type
	// for any unrecognized message types, so we just accept whatever comes in
	// and store the body, as is.
	// If the caller wants to verify the type, they should do so before calling this method.
	// We only accept the body and routingkey, all other fields are ignored/lost.
	u.RkFromRabbitMQ = m.RoutingKey
	u.Body = m.Body
	return nil
}

type RMQIdentifier interface {
	RMQRoutingIdentifier() string
	RMQDataIdentifier() string
}

type RMQDataStorage interface {
	RMQDataStoreTable() string
	RMQDataStoreIndex() string
}

type RMQSerializationOnWire interface {
	RMQEncodeMessage() (*MessageOverRabbitMQ, error)
	RMQDecodeMessage(m *MessageOverRabbitMQ) error
}

var _ = []interface {
	RMQIdentifier
	RMQSerializationOnWire
}{
	(*Unknown)(nil),
}
