package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var LogLevel = map[string]int{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
}

var timestamp = time.Now().UTC().Format(time.DateTime)

type Logger struct {
	Level int
}

func New(level string) *Logger {
	levelFromMap, ok := LogLevel[strings.ToUpper(level)]
	if !ok {
		levelFromMap = ERROR
	}
	return &Logger{
		Level: levelFromMap,
	}
}

func (l Logger) Error(msg string) {
	if l.Level <= ERROR {
		fmt.Fprintf(os.Stdout, "[ERROR] %s %s\n", timestamp, msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.Level <= DEBUG {
		fmt.Fprintf(os.Stdout, "[DEBUG] %s %s\n", timestamp, msg)
	}
}

func (l Logger) Info(msg string) {
	if l.Level <= INFO {
		fmt.Fprintf(os.Stdout, "[INFO] %s %s\n", timestamp, msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.Level <= WARN {
		fmt.Fprintf(os.Stdout, "[WARN] %s %s\n", timestamp, msg)
	}
}
