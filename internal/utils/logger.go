package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.WithTime(time.Now())
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	return logger
}
