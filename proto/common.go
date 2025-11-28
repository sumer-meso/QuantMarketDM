package proto

import (
	"context"
	"fmt"

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

func (m *MessageOverRabbitMQ) PublishInfoOnWire() amqp.Publishing {
	return amqp.Publishing{
		ContentType: "application/json",
		Headers: amqp.Table{
			HeaderKeyType:  m.DataIdentifier,
			HeaderKeyTable: m.StoreTable,
			HeaderKeyIndex: m.StoreIndex,
		},
		Expiration:   defaultMsgTTLInQueue,
		DeliveryMode: amqp.Transient,
		Body:         m.Body,
	}
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
	// The only reason we have this func here is to be used as we coverting to Known type.
	// In normal usage, Unknown type should not be encoded back to RMQ message.
	// since we keep using the body bytes pointer as is, so not much perf lost here.
	return &MessageOverRabbitMQ{
		RoutingKey:     u.RkFromRabbitMQ,
		DataIdentifier: u.RMQDataIdentifier(),
		StoreTable:     "N/A",
		StoreIndex:     "N/A",
		Body:           u.Body,
	}, nil
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

func Ptr2Str[T any](v *T) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", *v)
}

type NotMatchError struct {
	Expected string // e.g., "orderbook"
	Actual   string // e.g., "trade"
}

func (e NotMatchError) Error() string {
	return fmt.Sprintf("data identifier not match, expected: %s, actual: %s", e.Expected, e.Actual)
}

// UnknownToKnown converts an Unknown type to a known type T.
// T must implement RMQSerializationOnWire and RMQIdentifier interfaces.
// If the DataIdentifier in Unknown does not match the expected type T,
// we forecefully set it to match T's DataIdentifier before decoding.
// This is not encouraged for general use, as it may lead to unexpected behavior
// if the underlying data does not actually match the expected type.
// Use with caution.
func UnknownToKnown[T any, PT interface {
	*T
	RMQSerializationOnWire
	RMQIdentifier
}](unk *Unknown) (*T, error) {
	var zero T
	if unk == nil {
		return nil, fmt.Errorf("[UnknownToKnown] input Unknown is nil")
	}

	v := PT(&zero)

	unkMsg, err := unk.RMQEncodeMessage()
	if err != nil {
		return nil, fmt.Errorf("[UnknownToKnown] failed to encode Unknown message: %w", err)
	}

	// Set the DataIdentifier to match the target type.
	// This ensures correct decoding, especially in polymorphic scenarios.
	// But this is not encouraged for general use.
	unkMsg.DataIdentifier = v.RMQDataIdentifier()
	if err := v.RMQDecodeMessage(unkMsg); err != nil {
		return nil, fmt.Errorf("[UnknownToKnown] failed to decode message: %w", err)
	}

	return &zero, nil
}
