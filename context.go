package log

import (
	"fmt"
	"time"
)

// Fields describes the format of named values used to populate the context of
// the log record.
type Fields map[string]interface{}

// Context describes context information for logging.
type Context struct {
	fields Fields
	logger *Logger
}

// AddField adds a new field to the context.
func (c *Context) AddField(name string, value interface{}) *Context {
	if c.fields == nil {
		c.fields = make(Fields, 1)
	}
	c.fields[name] = value
	return c
}

// // AddFields adds new fields to the context.
func (c *Context) AddFields(fields Fields) *Context {
	if c.fields == nil {
		c.fields = fields
	} else {
		for name, value := range fields {
			c.fields[name] = value
		}
	}
	return c
}

// WithFields creates a new context for logging, adding a new field list.
func (c *Context) WithFields(fields Fields) *Context {
	for name, value := range c.fields {
		if _, ok := fields[name]; !ok {
			fields[name] = value
		}
	}
	return &Context{fields: fields, logger: c.logger}
}

// WithField creates a new context for logging, adding the new named field.
func (c *Context) WithField(name string, value interface{}) *Context {
	fields := make(Fields, len(c.fields)+1)
	fields[name] = value
	return c.WithFields(fields)
}

// WithError creates a new context for logging, adding the error field.
func (c *Context) WithError(err error) *Context {
	return c.WithField("error", err)
}

// WithSource creates a new context for logging, adding the caller source field.
func (c *Context) WithSource(calldepth int) *Context {
	return c.WithField("source", MakeCaller(calldepth+1))
}

func (c *Context) print(level Level, message string) {
	entry := entries.Get().(*Entry)
	entry.Timestamp = time.Now()
	entry.Level = level
	entry.Message = message
	entry.Fields = c.fields
	entry.Source = nil
	c.logger.handle(entry)
	entries.Put(entry)
}

// Debug displays the debug message in the log.
func (c *Context) Debug(message string) {
	c.print(DebugLevel, message)
}

// Debugf displays the debug formatted message in the log.
func (c *Context) Debugf(format string, v ...interface{}) {
	c.print(DebugLevel, fmt.Sprintf(format, v...))
}

// Info displays the message in the log.
func (c *Context) Info(message string) {
	c.print(InfoLevel, message)
}

// Infof displays the formatted message in the log.
func (c *Context) Infof(format string, v ...interface{}) {
	c.print(InfoLevel, fmt.Sprintf(format, v...))
}

// Error displays the error message in the log.
func (c *Context) Error(message string) {
	c.print(ErrorLevel, message)
}

// Errorf displays the formatted error message in the log.
func (c *Context) Errorf(format string, v ...interface{}) {
	c.print(ErrorLevel, fmt.Sprintf(format, v...))
}
