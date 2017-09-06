package log

import (
	"fmt"
	"time"
)

// TraceContext describes the context for the trace. Returned by the call to
// the Trace method.
type TraceContext struct {
	Message string
	context *Context
	started time.Time
}

// Stop should be used with Trace, to fire off the completion message. When an
// err is passed the "error" field is set, and the log level is error.
func (t *TraceContext) Stop(err *error) error {
	t.WithField("duration", time.Since(t.started))
	level := DebugLevel
	if err != nil && *err != nil {
		t.WithField("error", (*err).Error())
		level = ErrorLevel
	}
	return t.context.print(level, t.Message)
}

// WithFields adds to the TraceContext additional fields.
func (t *TraceContext) WithFields(fields Fields) *TraceContext {
	t.context.Fields = t.context.Fields.WithFields(fields)
	return t
}

// WithField adds to the TraceContext additional named field.
func (t *TraceContext) WithField(name string, value interface{}) *TraceContext {
	t.context.Fields = t.context.Fields.WithField(name, value)
	return t
}

func (c *Context) trace(message string) *TraceContext {
	return &TraceContext{
		Message: message,
		context: c.newContext(c.Fields),
		started: time.Now(),
	}
}

// Trace sends to the log debug message and returns a new TraceContext
// with a Stop method to fire off a corresponding completion log. Useful with
// defer.
func (c *Context) Trace(message string) *TraceContext {
	// do not move to the trace method to operate correctly determining Source
	// c.print(DebugLevel, message)
	return c.trace(message)
}

// Tracef sends to the log formatted debug message and returns a new
// TraceContext with a Stop method to fire off a corresponding completion log.
// Useful with defer.
func (c *Context) Tracef(format string, v ...interface{}) *TraceContext {
	message := fmt.Sprintf(format, v...)
	// c.print(DebugLevel, message)
	return c.trace(message)
}
