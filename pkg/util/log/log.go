package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

type Logger struct {
	out io.Writer
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{out}
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(LevelInfo, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(LevelWarn, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log(LevelError, format, v...)
}

func (l *Logger) log(level string, format string, v ...interface{}) {
	timestamp := time.Now().Format(time.RFC3339)
	msg := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, fmt.Sprintf(format, v...))
	fmt.Fprint(l.out, msg)
}

var defaultLogger = NewLogger(os.Stdout)

func Infof(format string, v ...interface{}) {
	defaultLogger.Info(format, v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warn(format, v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Error(format, v...)
}
