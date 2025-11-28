package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ctx      context.Context
	client   *Client
	Exchange string
	channel  *amqp.Channel
}

func (c *Client) NewPublisher(ctx context.Context, exchange string) *Publisher {
	pub := &Publisher{
		ctx:      ctx,
		client:   c,
		Exchange: exchange,
	}
	c.publishers = append(c.publishers, pub)
	go pub.recover()
	return pub
}

func (p *Publisher) recover() {
	p.client.wgGlobal.Add(1)
	defer p.client.wgGlobal.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		conn := p.client.getConn()
		if conn == nil {
			timeSleepSelect(p.ctx, 100)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			timeSleepSelect(p.ctx, 100)
			continue
		}

		if p.Exchange != "" {
			ch.ExchangeDeclare(p.Exchange, "topic", true, false, false, false, nil)
		}

		p.channel = ch
		return
	}
}

func (p *Publisher) Publish(rk string, info amqp.Publishing) error {
	for {
		select {
		case <-p.ctx.Done():
			return context.Canceled
		default:
		}

		if p.channel == nil {
			timeSleepSelect(p.ctx, 100)
			continue
		}

		err := p.channel.PublishWithContext(
			p.ctx,
			p.Exchange,
			rk,
			false, false,
			info,
		)

		if err != nil {
			timeSleepSelect(p.ctx, 200)
			continue
		}

		return nil
	}
}
