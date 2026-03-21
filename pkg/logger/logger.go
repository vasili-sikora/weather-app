package logger

import (
	"log"
	"os"
)

type StdLogger struct {
	infoLogger  *log.Logger
	debugLogger *log.Logger
	errorLogger *log.Logger
	debugMode   bool
}

func New(debugMode bool) *StdLogger {
	return &StdLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		debugMode:   debugMode,
	}
}

func (l *StdLogger) Info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *StdLogger) Debug(msg string) {
	if l.debugMode {
		l.debugLogger.Println(msg)
	}
}

func (l *StdLogger) Error(msg string) {
	l.errorLogger.Println(msg)
}
