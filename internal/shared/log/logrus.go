package log

import "github.com/sirupsen/logrus"

func NewLogrus() Logger {
	return logrus.New()
}
