package log

import (
	"bytes"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"sort"
	"strings"
	"sync"
)

// DefaultPadding длина сообщения для выравнивания по умолчанию.
var DefaultPadding = 32

// ConsoleHandler описывает обработчик для консольного лог, который является
// "обвязкой" над  стандартной библиотекой log.
type ConsoleHandler struct {
	Prefixes       map[Level]string // префиксы для вывода сообщений в лог
	Padding        int              // используется для выравнивая сообщений
	*stdlog.Logger                  // стандартный лог
}

// New возвращает новый инициализированный обработчик консольного
// лога. Если в параметре передано значение nil, то будет создан новый лог по
// умолчанию: с выводом в stderr и стандартным набором флагов для атрибутов.
func NewConsoleHandler(w io.Writer, flags int) *ConsoleHandler {
	// если параметры лога не указаны, то создаем стандартный лог
	if w == nil {
		w = os.Stderr
	}
	// инициализируем консольный лог
	return &ConsoleHandler{
		Logger:  stdlog.New(w, "", flags),
		Padding: DefaultPadding,
	}
}

// Handle обеспечивает поддержку вывода записей в лог.
func (h *ConsoleHandler) Handle(e *Entry) error {
	buf := buffers.Get().(*bytes.Buffer) // получаем буфер из пула
	buf.Reset()
	// если заданы префиксы для лога, то добавляем их к сообщению
	if prefix := h.Prefixes[e.Level]; prefix != "" {
		fmt.Fprintf(buf, "%s ", prefix)
	}
	// записываем само сообщение
	buf.WriteString(e.Message)
	// если есть поля сообщения, то добавляем их в виде строки
	if len(e.Fields) > 0 {
		// добавляем пробелы для выравнивания
		if length := len(e.Message); length < h.Padding {
			buf.WriteString(strings.Repeat(" ", h.Padding-length))
		}
		// выбираем названия полей для сортировки
		names := make([]string, 0, len(e.Fields))
		for name := range e.Fields {
			names = append(names, name)
		}
		sort.Strings(names)
		// выводим поля в отсортированном виде
		for _, name := range names {
			fmt.Fprintf(buf, " %s=", name) // записываем имя параметра
			// приводим значение к строковому виду
			value := fmt.Sprint(e.Fields[name])
			// если оно содержит пробелы, то выводим его в кавычках
			if strings.ContainsAny(value, " \t\n\r\"'") {
				fmt.Fprintf(buf, "%q", value)
			} else {
				buf.WriteString(value)
			}
		}
	}
	h.Logger.Output(5, buf.String()) // выводим в лог
	buffers.Put(buf)                 // освобождаем буфер
	return nil
}

// buffers содержит пул буферов для формирования сообщений.
var buffers = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
