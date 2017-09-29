package log

import (
	"fmt"
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
			buf = append(buf, "96"...)
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
		switch value := field.Value.(type) {
		case nil:
			buf = append(buf, "nil"...)
		case string:
			buf = append(buf, value...)
		case error:
			buf = strconv.AppendQuote(buf, value.Error())
			buf = append(buf, " \x1b[2m[\x1b[0m\x1b[91m"...)
			buf = append(buf, reflect.TypeOf(value).String()...)
			buf = append(buf, "\x1b[0m\x1b[2m]\x1b[0m"...)
		case bool:
			buf = strconv.AppendBool(buf, value)
		case int:
			buf = strconv.AppendInt(buf, int64(value), 10)
		case int8:
			buf = strconv.AppendInt(buf, int64(value), 10)
		case int16:
			buf = strconv.AppendInt(buf, int64(value), 10)
		case int32:
			buf = strconv.AppendInt(buf, int64(value), 10)
		case int64:
			buf = strconv.AppendInt(buf, value, 10)
		case uint:
			buf = strconv.AppendUint(buf, uint64(value), 10)
		case uint8:
			buf = strconv.AppendUint(buf, uint64(value), 10)
		case uint16:
			buf = strconv.AppendUint(buf, uint64(value), 10)
		case uint32:
			buf = strconv.AppendUint(buf, uint64(value), 10)
		case uint64:
			buf = strconv.AppendUint(buf, value, 10)
		case float32:
			buf = strconv.AppendFloat(buf, float64(value), 'g', -1, 32)
		case float64:
			buf = strconv.AppendFloat(buf, value, 'g', -1, 64)
		case time.Time:
			buf = append(buf, '"')
			if !value.IsZero() {
				buf = value.AppendFormat(buf, "2006-01-02 15:04:05")
			}
			buf = append(buf, '"')
		case fmt.Stringer:
			buf = append(buf, value.String()...)
		default:
			buf = append(buf, fmt.Sprint(value)...)
		}
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
