package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/AmFlint/taco-api-go/constants"
)

func GenerateLogger(handlerType, path, method string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		constants.HandlerKeyLogger: handlerType,
		constants.PathKeyLogger: path,
		constants.HttpVerbKeyLogger: method,
	})
}
