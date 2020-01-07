package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "timestamp",
			log.FieldKeyLevel: "log_level",
			log.FieldKeyMsg:   "message",
		},
		TimestampFormat: "2006-01-02T15:04:05.000-07:00",
	})
	log.SetOutput(os.Stdout)

	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if level, err := log.ParseLevel(v); err != nil {
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetLevel(level)
		}
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
