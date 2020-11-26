package config

import (
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

func InitializeLogger(logFile string) {
	log.SetLevel(getLogLevel())
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		PadLevelText:    true,
		TimestampFormat: time.RFC1123,
	})

	multiWriter := io.MultiWriter(colorable.NewColorableStdout(), &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	log.SetOutput(multiWriter)
}

func getLogLevel() log.Level {
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "DEBUG":
		return log.DebugLevel
	case "INFO":
		return log.InfoLevel
	case "TRACE":
		return log.TraceLevel
	default:
		return log.ErrorLevel
	}
}
