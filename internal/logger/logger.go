package logger

import (
	"io"
	"log"
	"os"
)

//counterfeiter:generate . Logger
type Logger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
	Fatal(message string)
	Close() error
}

type logger struct {
	logFile *os.File
	logger  *log.Logger
}

func NewLogger(logFilePath string) Logger {
	var writers []io.Writer
	writers = append(writers, os.Stdout) // Always write to console

	var logFile *os.File
	if logFilePath != "" {
		var err error
		logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Printf("Failed to open log file %s: %v", logFilePath, err)
		} else {
			writers = append(writers, logFile)
		}
	}

	multiWriter := io.MultiWriter(writers...)
	loggerInstance := log.New(multiWriter, "", log.LstdFlags)

	return &logger{
		logFile: logFile,
		logger:  loggerInstance,
	}
}

func (l *logger) Info(message string) {
	l.logger.Println("[INFO] " + message)
}

func (l *logger) Warn(message string) {
	l.logger.Println("[WARN] " + message)
}

func (l *logger) Error(message string) {
	l.logger.Println("[ERROR] " + message)
}

func (l *logger) Fatal(message string) {
	l.logger.Println("[FATAL] " + message)
	os.Exit(1)
}

func (l *logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
