package rabbitmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sumer-meso/QuantMarketDM/common"
	"github.com/sumer-meso/QuantMarketDM/utils/logging"
)

type Client struct {
	ctx        context.Context
	conn       *amqp.Connection
	lock       sync.RWMutex
	cfg        common.RabbitMQ
	consumers  []*Consumer
	publishers []*Publisher

	queues   map[string][]string
	wgGlobal sync.WaitGroup
}

func NewClient(ctx context.Context, cfg common.RabbitMQ) *Client {
	c := &Client{
		cfg:        cfg,
		ctx:        ctx,
		queues:     make(map[string][]string),
		consumers:  make([]*Consumer, 0),
		publishers: make([]*Publisher, 0),
	}
	go c.reconnectLoop()
	return c
}

// WaitAndClose waits for all goroutines to finish and closes the connection.
func (c *Client) WaitAndClose() {
	c.wgGlobal.Wait()
	c.lock.Lock()
	if c.conn != nil {
		c.conn.Close()
	}
	c.lock.Unlock()
}

func (c *Client) rmqConnectStr() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", c.cfg.Username, c.cfg.Password, c.cfg.Host, c.cfg.Port)
}

func (c *Client) reconnectLoop() {
	backoff := time.Second
	c.wgGlobal.Add(1)
	defer c.wgGlobal.Done()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		conn, err := amqp.Dial(c.rmqConnectStr())
		if err != nil {
			logging.Logf("[RabbitMq][Client] connect failed: %v", err)
			time.Sleep(backoff)
			backoff = min(backoff*2, 30*time.Second)
			continue
		}

		backoff = time.Second

		c.lock.Lock()
		c.conn = conn
		c.lock.Unlock()

		logging.Logf("[RabbitMq][Client] Connected")

		c.recoverAll()

		notify := conn.NotifyClose(make(chan *amqp.Error))
		select {
		case closeErr := <-notify:
			logging.Logf("[AMQP] Connection lost: %v", closeErr)
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Client) recoverAll() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, cons := range c.consumers {
		go cons.recover()
	}
	for _, pub := range c.publishers {
		go pub.recover()
	}
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func (c *Client) getConn() *amqp.Connection {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.conn
}
