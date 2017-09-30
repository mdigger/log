package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Encoder описывает интерфейс для форматирования записей лога. Используется
// Writer для задания формата. Данная библиотека содержит поддержку двух
// форматов логов: Console и JSON.
type Encoder interface {
	Encode(entry *Entry) []byte
}

// Writer описывает обработчик лога, записывающего в файл, консоль или
// другой поток.
type Writer struct {
	enc Encoder
	lvl Level
	w   io.Writer
	mu  sync.RWMutex
	Logger
}

// NewWriter возвращает новый обработчик лога.
func NewWriter(w io.Writer, lvl Level, enc Encoder) *Writer {
	if enc == nil {
		enc = new(Console)
	}
	var h = &Writer{w: w, lvl: lvl, enc: enc}
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
func (h *Writer) SetFormat(enc Encoder) {
	if enc == nil {
		enc = new(Console)
	}
	h.mu.Lock()
	h.enc = enc
	h.mu.Unlock()
}

// String возвращает уровень и формат вывода лога.
func (h *Writer) String() string {
	h.mu.RLock()
	var level string
	switch h.lvl {
	case -128:
		level = "ALL"
	case 127:
		level = "NONE"
	case TRACE:
		level = "TRACE"
	case DEBUG:
		level = "DEBUG"
	case INFO:
		level = "INFO"
	case ERROR:
		level = "ERROR"
	case FATAL:
		level = "FATAL"
	default:
		level = strconv.Itoa(int(h.lvl))
	}
	switch h.enc.(type) {
	case *JSON:
		level += ":JSON"
	case *Color:
		level += ":DEV"
	case *Console:
	}
	h.mu.RUnlock()
	return level
}

// Set устанавливает уровень и формат вывода лога.
func (h *Writer) Set(opt string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, opt := range strings.FieldsFunc(opt, func(c rune) bool {
		return map[rune]bool{
			' ': true, '\t': true, ';': true, ',': true, ':': true}[c]
	}) {
		switch strings.ToUpper(opt) {
		case "ALL", "A", "*":
			h.lvl = -128
		case "TRACE", "TRC", "T":
			h.lvl = TRACE
		case "DEBUG", "DBG", "D":
			h.lvl = DEBUG
		case "INFO", "INF", "I":
			h.lvl = INFO
		case "WARNING", "WARN", "WRN", "W":
			h.lvl = WARN
		case "ERROR", "ERR", "E":
			h.lvl = ERROR
		case "FATAL", "FTL", "F":
			h.lvl = FATAL
		case "NONE", "NO", "N", "OFF", "FALSE":
			h.lvl = 127
		case "JSON", "JSN", "J":
			h.enc = new(JSON)
		case "STANDART", "STD", "S", "CONSOLE":
			h.enc = &Console{TimeFormat: "2006-01-02 15:04:05"}
		case "COLORS", "COLOR", "DEVELOPERS", "DEVELOPER", "DEVELOP", "DEV", "C":
			h.enc = &Color{KeyIndent: 8}
		case "":
		default:
			if lvl, err := strconv.ParseInt(opt, 10, 8); err == nil {
				h.lvl = Level(lvl)
				continue
			}
			return fmt.Errorf("unknown log format %q", opt)
		}
	}
	return nil
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
func (h *Writer) Write(lvl Level, category, msg string, fields []Field) error {
	h.mu.RLock()
	if h.enc == nil || h.w == nil || lvl < h.lvl {
		h.mu.RUnlock()
		return nil
	}
	h.mu.RUnlock()
	var entry = NewEntry(lvl, category, msg, fields)
	var buf = h.enc.Encode(entry)
	entry.Free()
	h.mu.Lock()
	_, err := h.w.Write(buf)
	h.mu.Unlock()
	buffers.Put(buf)
	return err
}
