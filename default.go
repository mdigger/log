package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	plainHandler = NewPlainHandler(os.Stdout, LstdFlags)
	std          = New(plainHandler)
)

// WithFields creates a new context for logging, adding a new field list.
func WithFields(fields Fields) *Context {
	return std.WithFields(fields)
}

// WithField creates a new context for logging, adding the new named field.
func WithField(name string, value interface{}) *Context {
	return std.WithField(name, value)
}

// WithError creates a new context for logging, adding the error field.
func WithError(err error) *Context {
	return std.WithError(err)
}

// WithError creates a new context for logging, adding the caller source field.
func WithSource(calldepth int) *Context {
	return std.WithSource(calldepth)
}

// Debug displays the debug message in the log.
func Debug(message string) {
	std.print(LevelDebug, message)
}

// Debugf displays the debug formatted message in the log.
func Debugf(format string, v ...interface{}) {
	std.print(LevelDebug, fmt.Sprintf(format, v...))
}

// Info displays the message in the log.
func Info(message string) {
	std.print(LevelInfo, message)
}

// Infof displays the formatted message in the log.
func Infof(format string, v ...interface{}) {
	std.print(LevelInfo, fmt.Sprintf(format, v...))
}

// Error displays the error message in the log.
func Error(message string) {
	std.print(LevelError, message)
}

// Error displays the formatted error message in the log.
func Errorf(format string, v ...interface{}) {
	std.print(LevelError, fmt.Sprintf(format, v...))
}

// Trace returns a new entry with a Stop method to fire off a corresponding
// completion log, useful with defer.
func Trace(message string) *Tracer {
	std.print(LevelInfo, message)
	return &Tracer{
		Message: message,
		context: std.Context,
		started: time.Now(),
	}
}

// Tracef outputs to the console information about the beginning of the trace
// and returns the trace context to further it is stopped by Stop method.
func Tracef(format string, v ...interface{}) *Tracer {
	message := fmt.Sprintf(format, v...)
	std.print(LevelInfo, message)
	return &Tracer{
		Message: message,
		context: std.Context,
		started: time.Now(),
	}
}

func Level() LogLevel {
	return plainHandler.Level()
}

// SetLevel sets the minimum event level that is supported by the logger.
func SetLevel(level LogLevel) {
	plainHandler.SetLevel(level)
}

// Flags returns the output flags for the logger.
func Flags() int {
	return plainHandler.Flags()
}

// SetFlags sets the output flags for the logger.
func SetFlags(flag int) {
	plainHandler.SetFlags(flag)
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	plainHandler.SetOutput(w)
}
