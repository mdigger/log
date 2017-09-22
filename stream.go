package log

import (
	"io"
	"os"
	"sync"
)

// StreamHandler описывает обработчик консольного лога.
type StreamHandler struct {
	format StreamFormatter
	lvl    Level
	w      io.Writer
	mu     sync.Mutex
}

// NewStreamHandler возвращает новый консольный обработчик лога.
//
// Если w не определен, то лог выводиться не будет. Все записи, с уровнем ниже
// указанного, будут игнорироваться. За представление формата записи отвечает
// format. Если он не задан, то используется консольный формат по умолчанию.
func NewStreamHandler(w io.Writer, lvl Level, format StreamFormatter) *StreamHandler {
	if format == nil {
		format = new(Console)
	}
	return &StreamHandler{w: w, lvl: lvl, format: format}
}

// Log отвечает за запись в лог.
func (h *StreamHandler) Log(lvl Level, category, msg string,
	fields ...interface{}) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.w == nil || lvl < h.lvl {
		return nil
	}
	return h.format.Log(h.w, lvl, category, msg, fields...)
}

// SetLevel устанавливает новый минимальный уровень для вывода в лог.
func (h *StreamHandler) SetLevel(lvl Level) {
	h.mu.Lock()
	h.lvl = lvl
	h.mu.Unlock()
}

// SetOutput переопределяет вывод лога. Если nil, то лог выводиться не будет.
func (h *StreamHandler) SetOutput(w io.Writer) {
	h.mu.Lock()
	h.w = w
	h.mu.Unlock()
}

// SetFormat задает свойства форматирования записей лога.
func (h *StreamHandler) SetFormat(format StreamFormatter) {
	if format == nil {
		format = new(Console)
	}
	h.mu.Lock()
	h.format = format
	h.mu.Unlock()
}

// IsTTY возвращает true, если поток является терминалом или файлом.
func (h *StreamHandler) IsTTY() bool {
	var tty bool
	if out, ok := h.w.(*os.File); ok {
		if fi, err := out.Stat(); err == nil {
			tty = fi.Mode()&(os.ModeDevice|os.ModeCharDevice) != 0
		}
		// tty = os.Getenv("TERM") != "dumb" &&
		// 	(isatty.IsTerminal(out.Fd()) ||
		// 		isatty.IsCygwinTerminal(out.Fd()))
	}
	return tty
}
