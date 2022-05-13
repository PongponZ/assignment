package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"assignment/model"
	"assignment/services"

	"github.com/streadway/amqp"
)

const (
	consumeAutoAck   = true
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
	consumerTag      = ""
)

type workerConsumer struct {
	amqpConn *amqp.Connection
}

func NewConsumer(amqpConn *amqp.Connection) *workerConsumer {
	return &workerConsumer{amqpConn: amqpConn}
}

func (c *workerConsumer) CreateChannel() (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *workerConsumer) StartConsumer(queueName string) error {

	ch, err := c.CreateChannel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	c.QueueDeclare(ch, queueName)

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	go c.Worker(deliveries)
	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))

	return chanErr

}

func (c *workerConsumer) Worker(messages <-chan amqp.Delivery) {

	for delivery := range messages {
		var source model.Source
		json.Unmarshal(delivery.Body, &source)

		fmt.Printf("%+v", source)

		result := services.GetIpInfo(source.IpAddress)

		dest := model.Destination{
			Username:  source.Username,
			IpAddress: source.IpAddress,
			Address: model.Address{
				CountryCode: result.CountryCode,
				CountryName: result.CountryName,
			},
		}

		publisher := NewPublish(c.amqpConn)
		publisher.Publish("destination", &dest)

	}
}

func (c *workerConsumer) QueueDeclare(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
