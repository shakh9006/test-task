package delivery

import (
	"github.com/shakh9006/rabbitmq-logger/utils"
	"github.com/streadway/amqp"
	"log"
)

const (
	exchangeLogger = "logger"
	exchangeAll    = "all_logs"
	errorQueueKey  = "error"
	infoQueueKey   = "info"
	debugQueueKey  = "debug"
	allQueueKey    = "all"
)

type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	handler    *ConsumerHandler
}

func NewRabbitMQ(connectionURI string, handler *ConsumerHandler) *RabbitMQ {
	conn, err := amqp.Dial(connectionURI)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	return &RabbitMQ{
		connection: conn,
		channel:    ch,
		handler:    handler,
	}
}

func (r *RabbitMQ) ExchangeDeclare(name string, kind string) {
	err := r.channel.ExchangeDeclare(
		name,
		kind,
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
		queType,
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare an error queue")

	// Bind the error queue to the exchange
	err = r.channel.QueueBind(
		q.Name,
		key,
		exchangeName,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to bind the error queue to the exchange")

	return q
}

func (r *RabbitMQ) GetConsumer(q amqp.Queue, key string) <-chan amqp.Delivery {
	data, err := r.channel.Consume(q.Name, key, false, false, false, false, nil)
	utils.FailOnError(err, "Failed to register a consumer")
	return data
}

func (r *RabbitMQ) Register() {
	defer r.connection.Close()
	defer r.channel.Close()

	r.ExchangeDeclare(exchangeLogger, "direct")
	r.ExchangeDeclare(exchangeAll, "fanout")

	errorQueue := r.QueueDeclare(exchangeLogger, errorQueueKey, "error")
	infoQueue := r.QueueDeclare(exchangeLogger, infoQueueKey, "info")
	debugQueue := r.QueueDeclare(exchangeLogger, debugQueueKey, "debug")
	allQueue := r.QueueDeclare(exchangeAll, allQueueKey, "")

	errorData := r.GetConsumer(errorQueue, "errorConsumer")
	infoData := r.GetConsumer(infoQueue, "infoConsumer")
	debugData := r.GetConsumer(debugQueue, "debugConsumer")
	allData := r.GetConsumer(allQueue, "allConsumer")

	forever := make(chan bool)

	r.handler.RunRoutine(errorData, "error")
	r.handler.RunRoutine(infoData, "info")
	r.handler.RunRoutine(debugData, "debug")
	r.handler.RunRoutine(allData, "all")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
