package rabbitMq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Rabbit struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
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

	err = amqpChannel.Qos(1, 0, false)
	if err != nil {
		log.Fatal("Could not declared QoS .")
	}

	return &Rabbit{
		Channel: amqpChannel,
		Queue:   queue,
	}
}
