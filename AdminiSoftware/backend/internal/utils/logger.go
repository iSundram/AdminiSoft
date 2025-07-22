
package utils

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *Logger) Error(message string) {
	l.errorLogger.Println(message)
}
package utils

import (
	"io"
	"log"
	"os"
	"time"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func NewFileLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	multiWriter := io.MultiWriter(os.Stdout, file)
	errorWriter := io.MultiWriter(os.Stderr, file)

	return &Logger{
		infoLogger:    log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warningLogger: log.New(multiWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger:   log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

func (l *Logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.warningLogger.Println(v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.warningLogger.Printf(format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

func (l *Logger) LogRequest(method, path, ip string, duration time.Duration, statusCode int) {
	l.infof("%s %s from %s - %d (%v)", method, path, ip, statusCode, duration)
}

func (l *Logger) LogError(err error, context string) {
	if err != nil {
		l.errorLogger.Printf("%s: %v", context, err)
	}
}

func (l *Logger) LogSecurityEvent(event, ip, details string) {
	l.warningLogger.Printf("Security Event - %s from %s: %s", event, ip, details)
}
package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *Logger) Error(msg string) {
	l.errorLogger.Println(msg)
}

func (l *Logger) Debug(msg string) {
	l.debugLogger.Println(msg)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}

func (l *Logger) LogRequest(method, path, ip string, duration time.Duration) {
	l.infof("%s %s from %s took %v", method, path, ip, duration)
}

func (l *Logger) infof(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}
