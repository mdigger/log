package log

import "sync"

// type Logger interface {
// 	WithFields(fields Fields) *Context
// 	WithField(name string, value interface{}) *Context
// 	WithError(err error) *Context
// 	WithSource(calldepth int) *Context
// 	Debug(message string)
// 	Debugf(format string, v ...interface{})
// 	Info(message string)
// 	Infof(format string, v ...interface{})
// 	Error(message string)
// 	Errorf(format string, v ...interface{})
// 	Trace(message string) *TraceContext
// 	Tracef(format string, v ...interface{}) *TraceContext
// }

// Logger represents a logger with configurable level and handlers.
type Logger struct {
	*Context           // empty log context
	handlers []Handler // handlers
	mu       sync.RWMutex
}

// New returns a new Logger is ready for logging. As parameters you can specify
// multiple handlers for the log.
func New(handlers ...Handler) *Logger {
	logger := &Logger{handlers: handlers}
	logger.Context = &Context{logger: logger}
	return logger
}

// AddHandler adds new handler for logs.
func (l *Logger) AddHandler(handlers ...Handler) {
	l.mu.Lock()
	l.handlers = append(l.handlers, handlers...)
	l.mu.Unlock()
}

func (l *Logger) handle(entry *Entry) (err error) {
	l.mu.RLock()
	for _, handler := range l.handlers {
		if herr := handler.Handle(entry); herr != nil {
			err = herr
		}
	}
	l.mu.RUnlock()
	return err
}
