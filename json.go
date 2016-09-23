package log

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"
)

// JSON represents a JSON logger handler with configurable Level.
type JSON struct {
	handler
}

// NewJSON returns a new initialized handler for the log in JSON format.
func NewJSON(w io.Writer, flag int) *JSON {
	var json = new(JSON)
	json.SetOutput(w)
	json.SetFlags(flag)
	return json
}

// Context returns a new Context for a JSON log.
func (h *JSON) Context() *Context {
	return NewContext(h, nil)
}

// Handle implements the Handler interface.
func (h *JSON) Handle(e *Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if e.Level < h.level || h.w == nil {
		return nil
	}
	flag := h.flag

	var jsonEntry = &struct {
		Timestamp string `json:"timestamp,omitempty"`
		*Entry
		Source string `json:"source,omitempty"`
	}{
		Entry: e,
	}

	timestamp := e.Timestamp
	if flag&LUTC != 0 {
		timestamp = timestamp.UTC()
	}
	switch flag & (Ldate | Ltime | Lmicroseconds) {
	case Ldate | Ltime | Lmicroseconds, Ldate | Lmicroseconds:
		jsonEntry.Timestamp = timestamp.Format(time.RFC3339Nano)
	case Ldate | Ltime:
		jsonEntry.Timestamp = timestamp.Format(time.RFC3339)
	case Ltime | Lmicroseconds, Lmicroseconds:
		jsonEntry.Timestamp = timestamp.Format("15:04:05.999999")
	case Ltime:
		jsonEntry.Timestamp = timestamp.Format("15:04:05")
	case Ldate:
		jsonEntry.Timestamp = timestamp.Format("2006-01-02")
	}

	if flag&(Llongfile|Lshortfile) != 0 {
		if e.Source == nil {
			e.Source = NewSource(5)
		}
		var file = e.Source.File
		if flag&Lshortfile != 0 {
			file = filepath.Base(file)
		}
		jsonEntry.Source = fmt.Sprintf("%s:%d", file, e.Source.Line)
	}

	enc := json.NewEncoder(h.w)
	if h.flag&Lindent != 0 {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(jsonEntry)
}
