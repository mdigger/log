package log

import (
	"fmt"
	"time"
)

// A Context represents an active logging object that generates lines of output
// to an Handler.
type Context struct {
	Fields  Fields
	handler Handler
}

// NewContext returns a new context initialized to specified processors with
// specified fields.
func NewContext(h Handler, fields Fields) *Context {
	return &Context{
		Fields:  fields,
		handler: h,
	}
}

func (c *Context) newContext(fields Fields) *Context {
	return NewContext(c.handler, fields)
}

// WithFields returns a new Context with added fields.
func (c *Context) WithFields(fields Fields) *Context {
	return c.newContext(c.Fields.WithFields(fields))
}

// WithField returns a new Context with added the named field.
func (c *Context) WithField(name string, value interface{}) *Context {
	return c.newContext(c.Fields.WithField(name, value))
}

// WithError returns a new Context with added field "error" contains the error
// description.
func (c *Context) WithError(err error) *Context {
	if err == nil {
		return c
	}
	return c.newContext(c.Fields.WithError(err))
}

// WithSource return new Context with added information about the file name and
// line number of the source code. Calldepth is the count of the number of
// frames to skip when computing the file name and line number. A value of 0
// will print the details for the caller.
func (c *Context) WithSource(calldepth int) *Context {
	return c.newContext(c.Fields.WithSource(calldepth + 1))
}

func (c *Context) print(level Level, message string) error {
	entry := entries.Get().(*Entry)
	entry.Timestamp = time.Now()
	entry.Level = level
	entry.Message = message
	entry.Fields = c.Fields
	entry.Source = nil
	err := c.handler.Handle(entry)
	entries.Put(entry)
	return err
}

// Info publishes the informational message to the log.
func (c *Context) Info(message string) error {
	return c.print(InfoLevel, message)
}

// Infof publishes the formatted informational message to the log.
func (c *Context) Infof(format string, v ...interface{}) error {
	return c.print(InfoLevel, fmt.Sprintf(format, v...))
}

// Debug publishes the debug message to the log.
func (c *Context) Debug(message string) error {
	return c.print(DebugLevel, message)
}

// Debugf publishes the formatted debug message to the log.
func (c *Context) Debugf(format string, v ...interface{}) error {
	return c.print(DebugLevel, fmt.Sprintf(format, v...))
}

// Warning publishes the warning message to the log.
func (c *Context) Warning(message string) error {
	return c.print(WarningLevel, message)
}

// Warningf publishes the formatted warning message to the log.
func (c *Context) Warningf(format string, v ...interface{}) error {
	return c.print(WarningLevel, fmt.Sprintf(format, v...))
}

// Error publishes the error message to the log.
func (c *Context) Error(message string) error {
	return c.print(ErrorLevel, message)
}

// Errorf publishes the formatted error message to the log.
func (c *Context) Errorf(format string, v ...interface{}) error {
	return c.print(ErrorLevel, fmt.Sprintf(format, v...))
}
