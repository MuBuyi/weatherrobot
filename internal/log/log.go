package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func Init() {
	Logger.Out = os.Stdout
	Logger.SetFormatter(&logrus.JSONFormatter{})
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}
