package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

const (
	exchangeName     = "logger"
	numberServiceURL = "http://numbers_service:8080/v1/number/"
	rabbitMQURL      = "amqp://guest:guest@rabbitmq:5672/"
)

type Record struct {
	RecordId string `json:"record_id"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func logMessage(ch *amqp.Channel, message string, severity string) {
	err := ch.Publish(
		exchangeName, // exchange
		severity,     // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s: %s\n", severity, message)
}

func main() {
	conn, err := amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the exchange
	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	// Declare the queue
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	// Bind the queue to the exchange with record_id key
	err = ch.QueueBind(
		q.Name,
		"record_id",
		exchangeName,
		false,
		nil,
	)
	failOnError(err, "Failed to bind the queue to the exchange")

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		//loggerConn, err := amqp.Dial(loggerStoreURL)
		//failOnError(err, "Error while connecting to RabbitMQ: %s")
		//defer loggerConn.Close()
		//
		//loggerCh, err := loggerConn.Channel()
		//failOnError(err, "Error while opening a channel: %s")
		//defer loggerCh.Close()

		for d := range msgs {
			var record Record
			err := json.Unmarshal(d.Body, &record)
			if err != nil {
				log.Printf("Error while unmarshalling message body: %s. Please sent data in {\"record_id\":\"123\"} format", err)
				continue
			}
			logMessage(ch, fmt.Sprintf("Received data with record_id: %s", record.RecordId), "debug")

			/**
			//defer resp.Body.Close()
			log.Printf("Response from user service: %s", resp)
			**/

			resp, err := http.Get(fmt.Sprintf("%s%s", numberServiceURL, record.RecordId))
			if err != nil {
				log.Printf("Error while making request to number service: %s", err)
				logMessage(ch, fmt.Sprintf("Error while making request to number service: %s", err), "error")
				continue
			}

			if resp.StatusCode == 500 {
				logMessage(ch, fmt.Sprintf("Request with request_id: %s got 500 status code", record.RecordId), "error")
			} else if resp.StatusCode == 400 {
				logMessage(ch, fmt.Sprintf("Request with request_id: %s not found", record.RecordId), "error")
			} else if resp.StatusCode == 200 {
				logMessage(ch, fmt.Sprintf("Request successfully received. request_id: %s logged to info", record.RecordId), "info")
			}

			log.Printf("message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
