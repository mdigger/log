package log

import (
	"reflect"
	"strconv"
	"time"
)

// Color выводит лог в более удобном для чтения в консоли виде, используя
// цветовые выделения и помещая параметры на новую строку.
type Color struct {
	Levels    map[Level]string // переопределение строк для вывода уровня
	KeyIndent int              // отступ от значения дополнительного параметра
}

// Encode форматирует в буфер запись лога для текстового консольного
// представления.
func (f Color) Encode(buf []byte, entry *Entry) []byte {
	// выводим время
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	buf = append(buf, "\x1b[2m"...)
	buf = entry.Timestamp.AppendFormat(buf, "15:04:05.000000")
	buf = append(buf, "\x1b[0m "...)
	// уровень записи
	level, ok := f.Levels[entry.Level]
	if !ok {
		switch entry.Level & -32 {
		case FATAL:
			level = "FATAL"
		case ERROR:
			level = "ERROR"
		case WARN:
			level = "WARN "
		case INFO:
			level = "INFO "
		case DEBUG:
			level = "DEBUG"
		case TRACE:
			level = "TRACE"
		default:
			level = ""
		}
	}
	if level != "" {
		buf = append(buf, "\x1b[7m\x1b["...)
		switch entry.Level & -32 {
		case FATAL:
			buf = append(buf, "35"...)
		case ERROR:
			buf = append(buf, "91"...)
		case WARN:
			buf = append(buf, "93"...)
		case INFO:
			buf = append(buf, "92"...)
		case DEBUG:
			buf = append(buf, "94"...)
		case TRACE:
			buf = append(buf, "95"...)
		default:
			buf = append(buf, "37"...)
		}
		buf = append(buf, 'm')
		buf = append(buf, level...)
		buf = append(buf, "\x1b[0m "...)
	}
	// категория
	if entry.Category != "" {
		buf = append(buf, "\x1b[2m[\x1b[0m\x1b[92m"...)
		buf = append(buf, entry.Category...)
		buf = append(buf, "\x1b[0m\x1b[2m]:\x1b[0m "...)
	}
	// основной текст
	if entry.Message != "" {
		buf = append(buf, entry.Message...)
	}
	// дополнительные поля
	for _, field := range entry.Fields {
		buf = append(buf, "\n    \x1b[36m"...)
		buf = append(buf, field.Name...)
		buf = append(buf, "\x1b[0m"...)
		for i := 0; i < f.KeyIndent-len(field.Name); i++ {
			buf = append(buf, ' ')
		}
		buf = append(buf, "\x1b[2m=\x1b[0m"...)
		if f.KeyIndent > 0 {
			buf = append(buf, ' ')
		}
		if e, ok := field.Value.(error); ok {
			buf = strconv.AppendQuote(buf, e.Error())
			buf = append(buf, " \x1b[2m[\x1b[0m\x1b[91m"...)
			buf = append(buf, reflect.TypeOf(e).String()...)
			buf = append(buf, "\x1b[0m\x1b[2m]\x1b[0m"...)
			continue
		}
		buf = consoleValue(buf, field.Value)
	}
	// для ошибок выводим стек вызовов
	if entry.Level >= WARN {
		if entry.Stack == nil {
			entry.CallStack(1)
		}
		for _, src := range entry.Stack {
			buf = append(buf, "\n  \x1b[2m- "...)
			buf = append(buf, src.Pkg...)
			buf = append(buf, "/\x1b[0m"...)
			buf = append(buf, src.File...)
			buf = append(buf, "\x1b[2m:\x1b[0m"...)
			buf = strconv.AppendInt(buf, int64(src.Line), 10)
			buf = append(buf, " \x1b[2m(\x1b[0m\x1b[36m"...)
			buf = append(buf, src.Func...)
			buf = append(buf, "\x1b[0m\x1b[2m)\x1b[0m"...)
		}
	}
	buf = append(buf, '\n')
	return buf
}
