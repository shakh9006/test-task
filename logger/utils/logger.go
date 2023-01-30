package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	log2 "log"
)

func init() {
	logLevel := GetEnvVar("LOG_LEVEL")

	if logLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	host, err := os.Hostname()
	if err != nil {
		log.Logger = log.With().Str("host", "unknown").Logger()
	} else {
		log.Logger = log.With().Str("host", host).Logger()
	}

	log.Logger = log.With().Str("service", "gin-boilerplate").Logger()

	log.Logger = log.With().Caller().Logger()
}

func FailOnError(err error, msg string) {
	if err != nil {
		log2.Panicf("%s: %s", msg, err)
	}
}
