package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	red    = 31
	yellow = 33
	blue   = 34
	cyan   = 36
)

type logFormatter struct {
	log.TextFormatter
}

func (f *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	var levelColor int

	switch entry.Level {
	case log.DebugLevel, log.TraceLevel:
		levelColor = cyan
	case log.WarnLevel:
		levelColor = yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = red
	case log.InfoLevel:
		levelColor = blue
	default:
		levelColor = blue
	}

	return []byte(
		fmt.Sprintf("[%s] - \x1b[%dm%s\x1b[0m - %s\n",
			entry.Time.Format("02-01-2006 15:04:05"),
			levelColor,
			strings.ToUpper(entry.Level.String()),
			entry.Message,
		)), nil
}

func initLogger() {
	log.SetFormatter(&logFormatter{})

	logLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		logLevel = log.InfoLevel
		log.Info("Defaulting loglevel to ", logLevel)
	}

	log.SetLevel(logLevel)
	log.Debug("Initialized logger with log level -> ", logLevel)
}
