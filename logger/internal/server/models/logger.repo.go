package models

import (
	"log"
	"os"
)

type LoggerRepository struct {
}

func (r *LoggerRepository) WriteLog(filename string, message string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(message)

	return nil
}

func NewLoggerRepository() *LoggerRepository {
	return &LoggerRepository{}
}
