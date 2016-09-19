package log

import (
	"sync"
	"time"
)

// Entry describes the entry for the output in the log.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Message   string    `json:"message"`
	Fields    Fields    `json:"fields,omitempty"`
	Source    *Caller   `json:"source,omitempty"`
}

var entries = sync.Pool{New: func() interface{} { return new(Entry) }}
