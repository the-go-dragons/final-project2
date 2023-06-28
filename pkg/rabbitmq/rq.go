package rabbitmq

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/the-go-dragons/final-project2/pkg/config"
)

var rabbitmqChannel *amqp.Channel

type SMSBody struct {
	Sender    string
	Receivers []string
	Massage   string
}

func Connect() {
	amqpServerURL := config.Config.Ribbitmq.Url

	// Create a new RabbitMQ connection.
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rabbitmqChannel, err = conn.Channel()
	if err != nil {
		panic(err)
	}
	defer rabbitmqChannel.Close()

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
}

func NewMassage(body SMSBody) {
	data, _ := json.Marshal(body)
	rabbitmqChannel.Publish("", "SMS", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(data),
	})
}
