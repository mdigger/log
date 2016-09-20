package log

import (
	"fmt"
	"time"
)

// TraceContext contains the execution context of a trace request.
type TraceContext struct {
	Message string
	context *Context
	started time.Time
}

// Stop stops the trace request and, depending on the error, generates a new
// entry in the log about the successful completion or error. Also in the log
// record is added to the duration of the query.
func (t *TraceContext) Stop(err *error) {
	context := t.context.WithField("duration", time.Since(t.started))
	if err == nil || *err == nil {
		context.print(InfoLevel, t.Message)
	} else {
		context.WithError(*err).print(ErrorLevel, t.Message)
	}
}

func (t *TraceContext) AddField(name string, value interface{}) *TraceContext {
	if t.context.fields == nil {
		t.context.fields = make(Fields, 1)
	}
	t.context.fields[name] = value
	return t
}

func (t *TraceContext) AddFields(fields Fields) *TraceContext {
	if t.context.fields == nil {
		t.context.fields = fields
	} else {
		for name, value := range fields {
			t.context.fields[name] = value
		}
	}
	return t
}

// Trace returns a new entry with a Stop method to fire off a corresponding
// completion log, useful with defer.
func (c *Context) Trace(message string) *TraceContext {
	c.print(InfoLevel, message)
	return &TraceContext{
		Message: message,
		context: c,
		started: time.Now(),
	}
}

// Tracef outputs to the console information about the beginning of the trace
// and returns the trace context to further it is stopped by Stop method.
func (c *Context) Tracef(format string, v ...interface{}) *TraceContext {
	message := fmt.Sprintf(format, v...)
	c.print(InfoLevel, message)
	return &TraceContext{
		Message: message,
		context: c,
		started: time.Now(),
	}
}
