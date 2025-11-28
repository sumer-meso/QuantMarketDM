package rabbitmq

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sumer-meso/QuantMarketDM/utils/logging"
)

const exchangeSuffix = "-market"
const messageTTLInQueue = 20000                // 20 seconds.
const queueTTLExpires = int32(120 * 60 * 1000) // 120 mintues.

func rk2Exchange(rk string) string {
	return strings.Split(rk, ".")[0] + exchangeSuffix
}

func (c *Client) DeclearQueueWithSomeBindings(name string, bindings []string) error {
	return c.declareQueue(name, bindings)
}

// DeclareQueue declares a queue with given name and bindings.
// If the queue already exists, it will update its bindings to match the given ones.
func (c *Client) declareQueue(name string, bindings []string) error {
	c.lock.Lock()
	c.queues[name] = bindings
	c.lock.Unlock()

	conn := c.getConn()
	if conn == nil {
		return fmt.Errorf("[RabbitMq][Queue] failed to get connection")
	}

	ch, err := conn.Channel()
	if err != nil {
		logging.Logf("[RabbitMq][Queue] Failed to get channel: %v", err)
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-message-ttl": messageTTLInQueue,
			"x-expires":     queueTTLExpires,
		},
	)
	if err != nil {
		logging.Logf("[RabbitMq][Queue] DeclareQueue failed: %v", err)
		return err
	}

	for _, rk := range bindings {
		exchange := rk2Exchange(rk)

		err = ch.QueueBind(
			name,
			rk,
			exchange,
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			logging.Logf("[RabbitMq][Queue] DeclareQueue QueueBind failed: %v", err)
			return err
		}
	}

	return c.deleteNonExistBindings(name, bindings)
}

// The input queue must exist.
func (c *Client) deleteNonExistBindings(queue string, bindings []string) error {
	if curBindings, err := c.ListBindings(queue); err != nil {
		logging.Logf("[RabbitMq][Queue] DeleteNonExistBindings ListBindings failed: %v", err)
		return err
	} else {
		for _, b := range curBindings {
			if !slices.Contains(bindings, b.RoutingKey) {
				c.queueUnbind(queue, b)
			}
		}
	}
	return nil
}

func (c *Client) queueUnbind(queue string, b Binding) error {
	if b.Source == "" {
		// We should ignore such bindiing,
		// since this looks like the default one created by system.
		return nil
	}
	conn := c.getConn()
	if conn == nil {
		return fmt.Errorf("[RabbitMq][Queue] failed to get connection")
	}

	ch, err := conn.Channel()
	if err != nil {
		logging.Logf("[RabbitMq][Queue] Failed to get channel: %v", err)
		return err
	}
	defer ch.Close()

	err = ch.QueueUnbind(
		queue,
		b.RoutingKey,
		b.Source,
		b.Arguments,
	)
	if err != nil {
		logging.Logf("[RabbitMq][Queue] Failed to unbind queue: %v", err)
		return err
	}

	logging.Logf("[RabbitMq][Queue] deleted a binding from queue %s: %v", queue, b)
	return nil
}

// Binding structure that returned from http management api.
type Binding struct {
	Source          string                 `json:"source"`
	Vhost           string                 `json:"vhost"`
	Destination     string                 `json:"destination"`
	DestinationType string                 `json:"destination_type"`
	RoutingKey      string                 `json:"routing_key"`
	PropertiesKey   string                 `json:"properties_key"`
	Arguments       map[string]interface{} `json:"arguments"`
}

// ListBindings from management api.
func (c *Client) ListBindings(queue string) ([]Binding, error) {
	endpoint := fmt.Sprintf("http://%s:%s/api/queues/%s/%s/bindings",
		c.cfg.Host,
		c.cfg.HttpPort,
		url.PathEscape("/"),
		url.PathEscape(queue),
	)
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("[RabbitMq][Queue] Failed to list bindings: %s", string(body))
	}

	var bindings []Binding
	if err := json.NewDecoder(resp.Body).Decode(&bindings); err != nil {
		return nil, err
	}
	return bindings, nil
}
