package log

import (
	"io"
	"sync"
)

// Handler is used to handle log events, outputting them to stdio or sending
// them to remote services.
type Handler interface {
	Handle(*Entry) error
}

func New(h ...Handler) *Context {
	return &Context{handler: handlers(h)}
}

type handlers []Handler

func (h handlers) Handle(e *Entry) (err error) {
	for _, handler := range h {
		if herr := handler.Handle(e); herr != nil {
			err = herr
		}
	}
	return err
}

type handler struct {
	level Level
	flag  int
	w     io.Writer
	mu    sync.Mutex
}

// Level returns the minimum event level that is supported by the logger.
func (h *handler) Level() Level {
	h.mu.Lock()
	level := h.level
	h.mu.Unlock()
	return level
}

// SetLevel sets the minimum event level that is supported by the logger.
func (h *handler) SetLevel(level Level) {
	h.mu.Lock()
	h.level = level
	h.mu.Unlock()
}

// Flags returns the output flags for the logger.
func (h *handler) Flags() int {
	h.mu.Lock()
	flag := h.flag
	h.mu.Unlock()
	return flag
}

// SetFlags sets the output flags for the logger.
func (h *handler) SetFlags(flag int) {
	h.mu.Lock()
	h.flag = flag
	h.mu.Unlock()
}

// SetOutput sets the output destination for the logger.
func (h *handler) SetOutput(w io.Writer) {
	h.mu.Lock()
	h.w = w
	h.mu.Unlock()
}
