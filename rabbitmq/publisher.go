package rabbitmq

import (
	"assignment/model"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

const (
	publishMandatory   = false
	publishImmediate   = false
	publishContentType = "application/octet-stream"
)

type workerPublisher struct {
	amqpConn *amqp.Connection
}

func NewPublish(amqpConn *amqp.Connection) *workerPublisher {
	return &workerPublisher{amqpConn: amqpConn}
}

func (p *workerPublisher) Publish(queueName string, body *model.Destination) {
	ch, err := p.CreateChannel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	p.QueueDeclare(ch, queueName)

	bodyJson, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish(
		"",
		queueName,
		publishMandatory,
		publishImmediate,
		amqp.Publishing{
			ContentType: publishContentType,
			Body:        bodyJson,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

}

func (p *workerPublisher) CreateChannel() (*amqp.Channel, error) {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (c *workerPublisher) QueueDeclare(ch *amqp.Channel, queueName string) error {
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
