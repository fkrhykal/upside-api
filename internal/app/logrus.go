package app

import (
	"os"

	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/sirupsen/logrus"
)

func NewLogrus() log.Logger {
	return &logrus.Logger{
		Out:       os.Stdout,
		Level:     logrus.DebugLevel,
		Hooks:     make(logrus.LevelHooks),
		Formatter: new(logrus.TextFormatter),
	}
}
