package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

var rabbitmqChannel *amqp.Channel

type SMSBody struct {
	Sender    string
	Receivers string
	Massage   string
}

func Connect() {
	if rabbitmqChannel != nil {
		return
	}
	amqpServerURL := config.Config.Ribbitmq.Url

	// Create a new RabbitMQ connection.
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}

	rabbitmqChannel, err = conn.Channel()
	if err != nil {
		panic(err)
	}

	_, err = rabbitmqChannel.QueueDeclare(
		"SMS", // queue name
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to RabbitMQ")
}

func NewMassage(body SMSBody) {
	data, _ := json.Marshal(body)
	err := rabbitmqChannel.Publish("", "SMS", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(data),
	})
	if err != nil {
		fmt.Println("Cant add queue: " + err.Error())
	} else {
		fmt.Println("Added to queue")
	}
}
