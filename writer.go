package log

import (
	"io"
	"os"
	"sync"
	"time"
)

// Formatter описывает интерфейс для форматирования записей лога. Используется
// Writer для задания формата. Данная библиотека содержит поддержку двух
// форматов логов: Console и JSON.
type Formatter interface {
	Format(buf []byte, entry *Entry) []byte
}

// Writer описывает обработчик лога, записывающего в файл, консоль или
// другой поток.
type Writer struct {
	format Formatter
	lvl    Level
	w      io.Writer
	mu     sync.RWMutex
	Logger
}

// NewWriter возвращает новый обработчик лога.
func NewWriter(w io.Writer, lvl Level, format Formatter) *Writer {
	if format == nil {
		format = new(Console)
	}
	var h = &Writer{w: w, lvl: lvl, format: format}
	h.Logger = Logger{h: h}
	return h
}

// SetLevel устанавливает новый минимальный уровень для вывода в лог.
func (h *Writer) SetLevel(lvl Level) {
	h.mu.Lock()
	h.lvl = lvl
	h.mu.Unlock()
}

// SetOutput переопределяет вывод лога. Если nil, то лог выводиться не будет.
func (h *Writer) SetOutput(w io.Writer) {
	h.mu.Lock()
	h.w = w
	h.mu.Unlock()
}

// SetFormat задает свойства форматирования записей лога.
func (h *Writer) SetFormat(format Formatter) {
	if format == nil {
		format = new(Console)
	}
	h.mu.Lock()
	h.format = format
	h.mu.Unlock()
}

// IsTTY возвращает true, если поток является терминалом или файлом.
func (h *Writer) IsTTY() bool {
	h.mu.RLock()
	if out, ok := h.w.(*os.File); ok {
		if fi, err := out.Stat(); err == nil {
			h.mu.RUnlock()
			return fi.Mode()&(os.ModeDevice|os.ModeCharDevice) != 0
		}
	}
	h.mu.RUnlock()
	return false
}

// Write поддерживает интерфейс записи логов Handler.
func (h *Writer) Write(lvl Level, calldepth int, category, msg string,
	fields []Field) error {
	h.mu.RLock()
	if h.format == nil || h.w == nil || lvl < h.lvl {
		h.mu.RUnlock()
		return nil
	}
	h.mu.RUnlock()
	var entry = entries.Get().(*Entry)
	entry.Timestamp = time.Time{} // не устанавливаем время до записи
	entry.Level = lvl
	entry.Category = category
	entry.Message = msg
	entry.Stack = nil // по умолчанию стек не заполняется
	entry.Fields = fields
	entry.calldepth = calldepth + 2
	var buf = buffers.Get().([]byte)
	buf = h.format.Format(buf[:0], entry)
	entries.Put(entry)
	h.mu.Lock()
	_, err := h.w.Write(buf)
	h.mu.Unlock()
	buffers.Put(buf)
	return err
}

var (
	buffers = sync.Pool{New: func() interface{} { return []byte{} }}
	entries = sync.Pool{New: func() interface{} { return new(Entry) }}
)
