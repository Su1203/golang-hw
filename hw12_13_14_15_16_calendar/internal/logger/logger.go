package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	level  LogLevel
	output io.Writer
}

func New(level string) *Logger {
	return &Logger{
		level:  parseLevel(level),
		output: os.Stdout,
	}
}

func parseLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN", "WARNING":
		return LevelWarn
	case "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

func (l *Logger) log(level LogLevel, levelStr, msg string) {
	if level < l.level {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(l.output, "%s [%s] %s\n", timestamp, levelStr, msg)
}

func (l *Logger) Debug(msg string) {
	l.log(LevelDebug, "DEBUG", msg)
}

func (l *Logger) Info(msg string) {
	l.log(LevelInfo, "INFO", msg)
}

func (l *Logger) Warn(msg string) {
	l.log(LevelWarn, "WARN", msg)
}

func (l *Logger) Error(msg string) {
	l.log(LevelError, "ERROR", msg)
}
