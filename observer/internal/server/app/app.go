package app

import (
	"github.com/shakh9006/rabbitmq-controller/internal/server/delivery"
	"github.com/shakh9006/rabbitmq-controller/utils"
)

func Run() {
	connectionURI := utils.GetEnvVar("RMQ_URL")
	mq := delivery.NewRabbitMQ(connectionURI)
	mq.Register()
}
