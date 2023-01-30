package delivery

import (
	"github.com/shakh9006/rabbitmq-logger/internal/server/services"
	"github.com/shakh9006/rabbitmq-logger/utils"
	"github.com/streadway/amqp"
	"log"
)

type ConsumerHandler struct {
	loggerService *services.LoggerService
}

func NewConsumerHandler(service *services.LoggerService) *ConsumerHandler {
	return &ConsumerHandler{
		loggerService: service,
	}
}

func (c *ConsumerHandler) RunRoutine(data <-chan amqp.Delivery, loggerType string) {
	go func() {
		for d := range data {
			log.Printf("Received %s; message: %s", loggerType, d.Body)
			err := c.loggerService.WriteLog(loggerType, d.Body)
			utils.FailOnError(err, "Failed to create a log row")
		}
	}()
}
