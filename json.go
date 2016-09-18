package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// JSONHandler описывает обработчик лога в формате JSON.
type JSONHandler struct {
	flag int // properties
	w    io.Writer
	mu   sync.Mutex
}

// New возвращает новый инициализированный обработчик лога в формате JSON.
func NewJSONHandler(w io.Writer, flag int) *JSONHandler {
	if w == nil {
		w = os.Stderr
	}
	return &JSONHandler{w: w, flag: flag}
}

func (h *JSONHandler) SetFlags(flag int) {
	h.mu.Lock()
	h.flag = flag
	h.mu.Unlock()
}

// Handle обеспечивает вывод записи в формате JSON.
func (h *JSONHandler) Handle(e *Entry) error {
	now := time.Now() // делаем это как можно раньше
	h.mu.Lock()
	// добавляем имя файла и номер строки исходного кода, если требуется
	if h.flag&(Llongfile|Lshortfile) != 0 {
		if _, file, line, ok := runtime.Caller(3); ok {
			if h.flag&Lshortfile != 0 {
				file = filepath.Base(file)
			}
			e.WithField("source", fmt.Sprintf("%s:%d", file, line))
		}
	}
	if h.flag&LUTC != 0 {
		now = now.UTC()
	}
	var timestamp string
	switch h.flag & (Ldate | Ltime | Lmicroseconds) {
	case Ldate | Ltime | Lmicroseconds, Ldate | Lmicroseconds:
		timestamp = time.RFC3339Nano
	case Ldate | Ltime:
		timestamp = time.RFC3339
	case Ltime | Lmicroseconds, Lmicroseconds:
		timestamp = "15:04:05.999999999"
	case Ltime:
		timestamp = "15:04:05" //Z07:00"
	case Ldate:
		timestamp = "2006-01-02"
	}
	// инициализируем кодировщик в формат JSON
	enc := json.NewEncoder(h.w)
	if h.flag&Lindent != 0 {
		enc.SetIndent("", "  ")
	}
	// выводим дату и время, если требуется
	var err error
	if timestamp != "" {
		timestamp = now.Format(timestamp)
		// записываем в лог
		err = enc.Encode(struct {
			Timestamp string `json:"timestamp"`
			*Entry
		}{timestamp, e})
	} else {
		err = enc.Encode(e) // не требуется даты и время
	}
	h.mu.Unlock()
	return err
}
