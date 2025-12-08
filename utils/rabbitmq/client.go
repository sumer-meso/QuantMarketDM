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

func NewClient(ctx context.Context, cfg common.RabbitMQ) (*Client, error) {
	c := &Client{
		cfg:        cfg,
		ctx:        ctx,
		queues:     make(map[string][]string),
		consumers:  make([]*Consumer, 0),
		publishers: make([]*Publisher, 0),
	}

	conn, err := c.connectWithRetry()
	if err != nil {
		return nil, err
	}
	c.setConn(conn)

	go c.reconnectLoop()
	return c, nil
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

func (c *Client) connectWithRetry() (*amqp.Connection, error) {
	backoff := time.Second

	for {
		select {
		case <-c.ctx.Done():
			return nil, c.ctx.Err()
		default:
		}

		conn, err := amqp.Dial(c.rmqConnectStr())
		if err == nil {
			logging.Logf("[RabbitMq][Client] connect OK")
			return conn, nil
		}

		logging.Logf("[RabbitMq][Client] connect failed: %v", err)

		select {
		case <-c.ctx.Done():
			return nil, c.ctx.Err()
		case <-time.After(backoff):
		}

		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func (c *Client) reconnectLoop() {
	c.wgGlobal.Add(1)
	defer c.wgGlobal.Done()

	for {

		conn := c.getConn()
		if conn == nil {
			newConn, err := c.connectWithRetry()
			if err != nil {
				return
			}
			c.setConn(newConn)
			c.recoverAll()
			continue
		}

		notify := conn.NotifyClose(make(chan *amqp.Error, 1))

		select {
		case <-c.ctx.Done():
			return
		case err, ok := <-notify:
			if !ok {
				logging.Logf("[RabbitMq][Client] connection closed (no error)")
			} else {
				logging.Logf("[RabbitMq][Client] connection lost: %v", err)
			}
		}

		newConn, err := c.connectWithRetry()
		if err != nil {
			return // ctx 已经被 cancel 之类的
		}
		c.setConn(newConn)
		c.recoverAll()
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

func (c *Client) setConn(conn *amqp.Connection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.conn = conn
}

func (c *Client) getConn() *amqp.Connection {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.conn
}
