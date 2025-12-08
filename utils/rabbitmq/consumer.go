package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sumer-meso/QuantMarketDM/utils/logging"
)

type Consumer struct {
	ctx         context.Context
	client      *Client
	Queue       string
	BufferSize  int
	WorkerCount int
	Handler     func(*amqp.Delivery)

	msgChan chan *amqp.Delivery
}

func (c *Client) NewConsumer(
	ctx context.Context, queue string,
	workerCount, buffer int,
	handler func(*amqp.Delivery),
) *Consumer {

	cons := &Consumer{
		ctx:         ctx,
		client:      c,
		Queue:       queue,
		WorkerCount: workerCount,
		BufferSize:  buffer,
		Handler:     handler,
		msgChan:     make(chan *amqp.Delivery, buffer),
	}

	c.consumers = append(c.consumers, cons)
	go cons.recover()

	return cons
}

func (c *Consumer) recover() {
	c.client.wgGlobal.Add(1)
	defer c.client.wgGlobal.Done()
	for {
		select {
		case <-c.ctx.Done():
			logging.Logf("[RabbitMq][Consumer] ctx cancelled outside. recover closing.")
			return
		default:
		}

		conn := c.client.getConn()
		if conn == nil {
			timeSleepSelect(c.ctx, 100)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			timeSleepSelect(c.ctx, 100)
			continue
		}

		ch.Qos(2000, 0, false)

		msgs, err := ch.Consume(c.Queue, "quant-market-rabbitmq", true, false, false, false, nil)
		if err != nil {
			timeSleepSelect(c.ctx, 100)
			continue
		}

		c.client.wgGlobal.Add(1)
		// read-loop
		go func() {
			defer c.client.wgGlobal.Done()
			for {
				select {
				case <-c.ctx.Done():
					logging.Logf("[RabbitMq][Consumer] ctx cancelled outside. read-loop closing.")
					return
				case m, ok := <-msgs:
					if !ok {
						return
					}
					c.msgChan <- &m
				}
			}
		}()

		c.client.wgGlobal.Add(c.WorkerCount)
		// worker pool
		for i := 0; i < c.WorkerCount; i++ {
			go func(idx int) {
				defer c.client.wgGlobal.Done()
				for {
					select {
					case <-c.ctx.Done():
						logging.Logf("[RabbitMq][Consumer] ctx cancelled outside. WorkerCount %d closing.", idx)
						return
					case msg := <-c.msgChan:
						c.Handler(msg)
					}
				}
			}(i)
		}

		return

	}

}
