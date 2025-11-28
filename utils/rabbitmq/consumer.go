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
	Handler     func(amqp.Delivery)

	msgChan chan amqp.Delivery
}

func (c *Client) NewConsumer(
	ctx context.Context, queue string,
	workerCount, buffer int,
	handler func(amqp.Delivery),
) *Consumer {

	cons := &Consumer{
		ctx:         ctx,
		client:      c,
		Queue:       queue,
		WorkerCount: workerCount,
		BufferSize:  buffer,
		Handler:     handler,
		msgChan:     make(chan amqp.Delivery, buffer),
	}

	c.consumers = append(c.consumers, cons)
	go cons.recover()

	return cons
}

func (c *Consumer) recover() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		conn := c.client.getConn()
		if conn == nil {
			timeSleepSelect(c.ctx, 500)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			timeSleepSelect(c.ctx, 500)
			continue
		}

		_, err = ch.QueueDeclare(c.Queue, true, false, false, false, nil)
		if err != nil {
			logging.Logf("[RabbitMq][Consumer] queue declare error: %v", err)
			return
		}

		ch.Qos(2000, 0, false)

		msgs, err := ch.Consume(c.Queue, "quant-market-rabbitmq", false, false, false, false, nil)
		if err != nil {
			timeSleepSelect(c.ctx, 500)
			continue
		}

		// read-loop
		go func() {
			for {
				select {
				case <-c.ctx.Done():
					return
				case m, ok := <-msgs:
					if !ok {
						return
					}
					c.msgChan <- m
				}
			}
		}()

		// worker pool
		for i := 0; i < c.WorkerCount; i++ {
			go func() {
				for {
					select {
					case <-c.ctx.Done():
						return
					case msg := <-c.msgChan:
						c.Handler(msg)
					}
				}
			}()
		}

		return

	}

}
