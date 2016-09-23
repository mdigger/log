package log

import (
	"sync"
	"time"
)

// Entry represents a single log entry.
//
// Source is used exclusively for caching the information about the source file
// and line number and is not automatically populated. If the log output is
// required this information, and this field is not filled in, it can be
// obtained by calling the NewSource function with the appropriate level of
// nesting of calls. It is desirable to store the information in the Source
// field, to avoid further calls to this expensive function by other possible
// processors of the log.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     Level     `json:"level,omitempty"`
	Message   string    `json:"message,omitempty"`
	Fields    Fields    `json:"fields,omitempty" bson:",omitempty"`
	Source    *Source   `json:"source,omitempty" bson:",omitempty"`
}

var entries = sync.Pool{New: func() interface{} { return new(Entry) }}
