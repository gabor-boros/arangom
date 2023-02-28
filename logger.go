package arangom

import (
	"fmt"
	"os"
)

var (
	// ErrNoLogger is returned when no logger is provided.
	ErrNoLogger = fmt.Errorf("no logger provided")
)

// Logger is an interface that is used to log messages.
type Logger interface {
	// Info logs an informational message.
	Info(args ...any)
	// Infof logs a formatted informational message.
	Infof(format string, args ...any)
	// Error logs an error message.
	Error(args ...any)
	// Errorf logs a formatted error message.
	Errorf(format string, args ...any)
	// Fatal logs a fatal message and exits.
	Fatal(args ...any)
	// Fatalf logs a formatted fatal message and exits.
	Fatalf(format string, args ...any)
}

// LogWriter is an interface that is used to write log messages.
type LogWriter interface {
	WriteString(s string) (n int, err error)
}

// DefaultLogger is a logger that prints messages to stdout and exits on fatal
// messages.
type DefaultLogger struct {
	Writer LogWriter
	Exiter func(int)
}

func (l *DefaultLogger) log(format, level string, args ...any) {
	_, _ = l.Writer.WriteString(fmt.Sprintf(format+"\n", level, fmt.Sprint(args...)))
	if level == "FATAL" {
		l.Exiter(1)
	}
}

func (l *DefaultLogger) Info(args ...any) {
	l.Infof("%s", fmt.Sprint(args...))
}

func (l *DefaultLogger) Infof(format string, args ...any) {
	l.log("[%s] %s", "INFO", fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Error(args ...any) {
	l.Errorf("%s", fmt.Sprint(args...))
}

func (l *DefaultLogger) Errorf(format string, args ...any) {
	l.log("[%s] %s", "ERROR", fmt.Sprintf(format, args...))
}

func (l *DefaultLogger) Fatal(args ...any) {
	l.Fatalf("%s", fmt.Sprint(args...))
}

func (l *DefaultLogger) Fatalf(format string, args ...any) {
	l.log("[%s] %s", "FATAL", fmt.Sprintf(format, args...))
}

// NewDefaultLogger returns a new DefaultLogger.
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		Writer: os.Stdout,
		Exiter: os.Exit,
	}
}
