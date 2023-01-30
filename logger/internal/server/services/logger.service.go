package services

import (
	"fmt"
	"github.com/shakh9006/rabbitmq-logger/internal/server/models"
)

type LoggerService struct {
	loggerRepository *models.LoggerRepository
}

func NewLoggerService(numberRepo *models.LoggerRepository) *LoggerService {
	return &LoggerService{
		loggerRepository: numberRepo,
	}
}

func (ls *LoggerService) WriteLog(logType string, message []byte) error {
	var logger models.Logger

	logger.Type = logType
	logger.Message = fmt.Sprintf("%s", message)

	filename := fmt.Sprintf("logs/%s.log", logType)
	fmt.Println("filename: ", filename)
	return ls.loggerRepository.WriteLog(filename, logger.Message)
}
