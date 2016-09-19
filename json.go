package log

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"time"
)

// JSONHandler describes the handler for logging to plain text.
type JSONHandler struct {
	w     io.Writer
	flag  int
	level LogLevel
	mu    sync.Mutex
}

// NewJSONHandler creates a new JSON logger Handler.
func NewJSONHandler(w io.Writer, flag int) *JSONHandler {
	return &JSONHandler{w: w, flag: flag}
}

// Level returns the minimum event level that is supported by the logger.
func (h *JSONHandler) Level() LogLevel {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.level
}

// SetLevel sets the minimum event level that is supported by the logger.
func (h *JSONHandler) SetLevel(level LogLevel) {
	h.mu.Lock()
	h.level = level
	h.mu.Unlock()
}

// Flags returns the output flags for the logger.
func (h *JSONHandler) Flags() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.flag
}

// SetFlags sets the output flags for the logger.
func (h *JSONHandler) SetFlags(flag int) {
	h.mu.Lock()
	h.flag = flag
	h.mu.Unlock()
}

// SetOutput sets the output destination for the logger.
func (h *JSONHandler) SetOutput(w io.Writer) {
	h.mu.Lock()
	h.w = w
	h.mu.Unlock()
}

// Handle implements Handler.
func (h *JSONHandler) Handle(e *Entry) error {
	h.mu.Lock()
	if e.Level < h.level {
		h.mu.Unlock()
		return nil
	}
	var jsonEntry = &struct {
		Timestamp string `json:"timestamp,omitempty"`
		*Entry
		Source string `json:"source,omitempty"`
	}{
		Entry: e,
	}

	timestamp := e.Timestamp
	if h.flag&LUTC != 0 {
		timestamp = timestamp.UTC()
	}
	switch h.flag & (Ldate | Ltime | Lmicroseconds) {
	case Ldate | Ltime | Lmicroseconds, Ldate | Lmicroseconds:
		jsonEntry.Timestamp = timestamp.Format(time.RFC3339Nano)
	case Ldate | Ltime:
		jsonEntry.Timestamp = timestamp.Format(time.RFC3339)
	case Ltime | Lmicroseconds, Lmicroseconds:
		jsonEntry.Timestamp = timestamp.Format("15:04:05.999999999")
	case Ltime:
		jsonEntry.Timestamp = timestamp.Format("15:04:05")
	case Ldate:
		jsonEntry.Timestamp = timestamp.Format("2006-01-02")
	}

	if h.flag&(Llongfile|Lshortfile) != 0 {
		if e.Source == nil {
			h.mu.Unlock()
			e.Source = getCaller(4)
			h.mu.Lock()
		}
		var file = e.Source.File
		if h.flag&Lshortfile != 0 {
			file = filepath.Base(file)
		}
		jsonEntry.Source = fmt.Sprintf("%s:%d", file, e.Source.Line)
	}

	enc := json.NewEncoder(h.w)
	if h.flag&Lindent != 0 {
		enc.SetIndent("", "  ")
	}
	err := enc.Encode(jsonEntry)
	h.mu.Unlock()
	return err
}
