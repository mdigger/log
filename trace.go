package log

import (
	"fmt"
	"time"
)

// Tracer contains the execution context of a trace request.
type Tracer struct {
	Message string
	context *Context
	started time.Time
}

// Stop stops the trace request and, depending on the error, generates a new
// entry in the log about the successful completion or error. Also in the log
// record is added to the duration of the query.
func (t *Tracer) Stop(err *error) {
	context := t.context.WithField("duration", time.Since(t.started))
	if err == nil || *err == nil {
		context.print(InfoLevel, t.Message)
	} else {
		context.WithError(*err).print(ErrorLevel, t.Message)
	}
}

// Trace returns a new entry with a Stop method to fire off a corresponding
// completion log, useful with defer.
func (c *Context) Trace(message string) *Tracer {
	c.print(InfoLevel, message)
	return &Tracer{
		Message: message,
		context: c,
		started: time.Now(),
	}
}

// Tracef outputs to the console information about the beginning of the trace
// and returns the trace context to further it is stopped by Stop method.
func (c *Context) Tracef(format string, v ...interface{}) *Tracer {
	message := fmt.Sprintf(format, v...)
	c.print(InfoLevel, message)
	return &Tracer{
		Message: message,
		context: c,
		started: time.Now(),
	}
}
