package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

var plainHandler = NewPlainHandler(os.Stdout, LstdFlags)

// GetLevel return current log level.
func GetLevel() Level {
	return plainHandler.Level()
}

// SetLevel sets the minimum event level that is supported by the logger.
func SetLevel(level Level) {
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

var log *logger

func init() {
	log = &logger{handlers: []Handler{plainHandler}}
	log.Context = &Context{logger: log}
}

func AddHandler(handlers ...Handler) {
	log.AddHandler(handlers...)
}

func Default() *Context {
	return log.Context
}

// WithFields creates a new context for logging, adding a new field list.
func WithFields(fields Fields) *Context {
	return log.WithFields(fields)
}

// WithField creates a new context for logging, adding the new named field.
func WithField(name string, value interface{}) *Context {
	return log.WithField(name, value)
}

// WithError creates a new context for logging, adding the error field.
func WithError(err error) *Context {
	return log.WithError(err)
}

// WithSource creates a new context for logging, adding the caller source field.
func WithSource(calldepth int) *Context {
	return log.WithSource(calldepth)
}

// Debug displays the debug message in the log.
func Debug(message string) {
	log.print(DebugLevel, message)
}

// Debugf displays the debug formatted message in the log.
func Debugf(format string, v ...interface{}) {
	log.print(DebugLevel, fmt.Sprintf(format, v...))
}

// Info displays the message in the log.
func Info(message string) {
	log.print(InfoLevel, message)
}

// Infof displays the formatted message in the log.
func Infof(format string, v ...interface{}) {
	log.print(InfoLevel, fmt.Sprintf(format, v...))
}

// Error displays the error message in the log.
func Error(message string) {
	log.print(ErrorLevel, message)
}

// Errorf displays the formatted error message in the log.
func Errorf(format string, v ...interface{}) {
	log.print(ErrorLevel, fmt.Sprintf(format, v...))
}

// Trace returns a new entry with a Stop method to fire off a corresponding
// completion log, useful with defer.
func Trace(message string) *TraceContext {
	log.print(InfoLevel, message)
	return &TraceContext{
		Message: message,
		context: log.Context,
		started: time.Now(),
	}
}

// Tracef outputs to the console information about the beginning of the trace
// and returns the trace context to further it is stopped by Stop method.
func Tracef(format string, v ...interface{}) *TraceContext {
	message := fmt.Sprintf(format, v...)
	log.print(InfoLevel, message)
	return &TraceContext{
		Message: message,
		context: log.Context,
		started: time.Now(),
	}
}
