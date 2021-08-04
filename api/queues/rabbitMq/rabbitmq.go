package rabbitMq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Rabbit struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func New(rabbitUrl, queueName string) *Rabbit {
	conn, err := amqp.Dial(rabbitUrl)

	if err != nil {
		log.Fatal("Can't connect to AMQP")
	}

	amqpChannel, err := conn.Channel()

	if err != nil {
		log.Fatal("Can't create a amqpChannel")
	}

	queue, err := amqpChannel.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		log.Fatal(fmt.Sprintf("Could not declared: %s queue.", queueName))
	}

	return &Rabbit{
		channel: amqpChannel,
		queue:   queue,
	}
}

func (r *Rabbit) PublishJob(body []byte) error {
	return r.channel.Publish("", r.queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
}
