package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/shakh9006/rabbitmq-controller/internal/server/models"
	"github.com/shakh9006/rabbitmq-controller/utils"
	"github.com/streadway/amqp"
	"io"
	"log"
	"net/http"
)

const (
	exchangeName     = "logger"
	numberServiceURL = "http://numbers_service:8080/v1/number/"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ(connectionURI string) *RabbitMQ {
	conn, err := amqp.Dial(connectionURI)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
	}
}

func (r *RabbitMQ) ExchangeDeclare(name string, kind string) {
	err := r.channel.ExchangeDeclare(
		name,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare an exchange")
}

func (r *RabbitMQ) QueueDeclare(exchangeName string, queType string, key string) amqp.Queue {
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare a queue")

	// Bind the queue to the exchange with record_id key
	err = r.channel.QueueBind(
		q.Name,
		"record_id",
		exchangeName,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to bind the queue to the exchange")

	return q
}

func (r *RabbitMQ) GetConsumer(q amqp.Queue) <-chan amqp.Delivery {
	data, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)
	utils.FailOnError(err, "Failed to register a consumer")

	return data
}

func (r *RabbitMQ) LogMessage(ch *amqp.Channel, message string, severity string) {
	err := ch.Publish(
		exchangeName, // exchange
		severity,     // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	utils.FailOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s: %s\n", severity, message)
}

func (r *RabbitMQ) Register() {
	defer r.connection.Close()
	defer r.channel.Close()

	r.ExchangeDeclare(exchangeName, "")
	errorQueue := r.QueueDeclare(exchangeName, "", "error")
	data := r.GetConsumer(errorQueue)

	forever := make(chan bool)

	go func() {
		for d := range data {
			var record models.Record
			err := json.Unmarshal(d.Body, &record)
			if err != nil {
				log.Printf("Error while unmarshalling message body: %s. Please sent data in {\"record_id\":\"123\"} format", err)
				continue
			}
			r.LogMessage(r.channel, fmt.Sprintf("Received data with record_id: %s", record.RecordId), "debug")

			resp, err := http.Get(fmt.Sprintf("%s%s", numberServiceURL, record.RecordId))

			if err != nil {
				log.Printf("Error while making request to number service: %s", err)
				r.LogMessage(r.channel, fmt.Sprintf("Error while making request to number service: %s", err), "error")
				continue
			}

			if resp.StatusCode == 500 {
				r.LogMessage(r.channel, fmt.Sprintf("Request with request_id: %s internal error", record.RecordId), "error")
			} else if resp.StatusCode == 404 {
				r.LogMessage(r.channel, fmt.Sprintf("Request with request_id: %s Not found", record.RecordId), "error")
			} else if resp.StatusCode == 200 {
				var response models.Response
				body, _ := io.ReadAll(resp.Body)
				if err := json.Unmarshal(body, &response); err != nil {
					log.Printf("Error while unmarshalling response body: %s", err)
					r.LogMessage(r.channel, fmt.Sprintf("Error while unmarshalling response body: %s", err), "error")
					continue
				}
				r.LogMessage(r.channel, fmt.Sprintf("Request with request_id: %s received.Number: %s", record.RecordId, response.Number), "info")
			}

			log.Printf("message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
