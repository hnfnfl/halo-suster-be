package logger

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	InfoLevel  = "info"
	ErrorLevel = "error"
	DebugLevel = "debug"
)

func NewLogger(logLevel string) (*logrus.Logger, error) {
	var level logrus.Level

	switch logLevel {
	case InfoLevel:
		level = logrus.InfoLevel
	case ErrorLevel:
		level = logrus.ErrorLevel
	case DebugLevel:
		level = logrus.DebugLevel
	default:
		return nil, fmt.Errorf("invalid log level: %s", logLevel)
	}

	formatter := &logrus.TextFormatter{
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	}
	logrus.SetFormatter(formatter)

	logger := &logrus.Logger{
		Out:          os.Stdout,
		Hooks:        nil,
		Formatter:    formatter,
		ReportCaller: false,
		Level:        level,
		ExitFunc:     nil,
	}

	return logger, nil
}

func ContextLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("LoggerForContext", log)
		c.Next()
	}
}

func FromContext(c *gin.Context) *logrus.Logger {
	logger, ok := c.Get("LoggerForContext")
	if !ok {
		fmt.Println("Logger not found in context")
		return nil
	}

	return logger.(*logrus.Logger)
}
