package log

import (
	"fmt"
	"io"
	"os"
)

var consoleHandler = NewConsole(os.Stderr, LstdFlags|Lindent)

// GetLevel return the default log level.
func GetLevel() Level {
	return consoleHandler.Level()
}

// SetLevel sets the minimum event level that is supported by the logger.
func SetLevel(level Level) {
	consoleHandler.SetLevel(level)
}

// Flags returns the output flags for the logger.
func Flags() int {
	return consoleHandler.Flags()
}

// SetFlags sets the output flags for the logger.
func SetFlags(flag int) {
	consoleHandler.SetFlags(flag)
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	consoleHandler.SetOutput(w)
}

var defaulContext = consoleHandler.Context()

func WithFields(fields Fields) *Context {
	return defaulContext.WithFields(fields)
}

func WithField(name string, value interface{}) *Context {
	return defaulContext.WithField(name, value)
}

func WithError(err error) *Context {
	return defaulContext.WithError(err)
}

// WithSource return new Context with added information about the file name and
// line number of the source code. Calldepth is the count of the number of
// frames to skip when computing the file name and line number. A value of 0
// will print the details for the caller.
func WithSource(calldepth int) *Context {
	return defaulContext.WithSource(calldepth + 1)
}

// Info publishes the informational message to the default log.
func Info(message string) error {
	return defaulContext.print(InfoLevel, message)
}

// Infof publishes the formatted informational message to the default log.
func Infof(format string, v ...interface{}) error {
	return defaulContext.print(InfoLevel, fmt.Sprintf(format, v...))
}

// Debug publishes the debug message to the default log.
func Debug(message string) error {
	return defaulContext.print(DebugLevel, message)
}

// Debugf publishes the formatted debug message to the default log.
func Debugf(format string, v ...interface{}) error {
	return defaulContext.print(DebugLevel, fmt.Sprintf(format, v...))
}

// Warning publishes the warning message to the default log.
func Warning(message string) error {
	return defaulContext.print(WarningLevel, message)
}

// Warningf publishes the formatted warning message to the default log.
func Warningf(format string, v ...interface{}) error {
	return defaulContext.print(WarningLevel, fmt.Sprintf(format, v...))
}

// Error publishes the error message to the default log.
func Error(message string) error {
	return defaulContext.print(ErrorLevel, message)
}

// Errorf publishes the formatted error message to the default log.
func Errorf(format string, v ...interface{}) error {
	return defaulContext.print(ErrorLevel, fmt.Sprintf(format, v...))
}

// Trace sends to the default log information message and returns a new
// TraceContext with a Stop method to fire off a corresponding completion log.
// Useful with defer.
func Trace(message string) *TraceContext {
	defaulContext.print(InfoLevel, message)
	return defaulContext.trace(message)
}

// Trace sends to the default log formatted information message and returns a
// new TraceContext with a Stop method to fire off a corresponding completion
// log. Useful with defer.
func Tracef(format string, v ...interface{}) *TraceContext {
	message := fmt.Sprintf(format, v...)
	defaulContext.print(InfoLevel, message)
	return defaulContext.trace(message)
}
