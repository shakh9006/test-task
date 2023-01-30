package app

import (
	"github.com/shakh9006/rabbitmq-logger/internal/server/delivery"
	"github.com/shakh9006/rabbitmq-logger/internal/server/models"
	"github.com/shakh9006/rabbitmq-logger/internal/server/services"
	"github.com/shakh9006/rabbitmq-logger/utils"
)

func Run() {
	connectionURI := utils.GetEnvVar("RMQ_URL")
	repo := models.NewLoggerRepository()
	service := services.NewLoggerService(repo)
	handlers := delivery.NewConsumerHandler(service)

	mq := delivery.NewRabbitMQ(connectionURI, handlers)
	mq.Register()
}
