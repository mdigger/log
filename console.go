package log

import (
	"fmt"
	"strconv"
	"time"
)

// Console поддерживает текстовое представление лога.
type Console struct {
	TimeFormat string           // формат вывода даты и времени
	UTC        bool             // вывод даты и времени в UTC
	WithSource bool             // выводить информацию об исходном коде
	Levels     map[Level]string // переопределение строк для вывода уровня
}

// Format форматирует в буфер запись лога для текстового консольного
// представления.
func (f Console) Format(buf []byte, entry *Entry) []byte {
	// выводим дату и время, если задан формат
	if f.TimeFormat != "" {
		if entry.Timestamp.IsZero() {
			entry.Timestamp = time.Now()
		}
		var ts = entry.Timestamp
		if f.UTC {
			ts = ts.UTC()
		}
		buf = ts.AppendFormat(buf, f.TimeFormat)
		buf = append(buf, ' ')
	}
	// информация об исходном файле
	if f.WithSource {
		if entry.Stack == nil {
			entry.CallStack(1)
		}
		if len(entry.Stack) > 0 {
			buf = append(buf, entry.Stack[0].Pkg...)
			buf = append(buf, '/')
			buf = append(buf, entry.Stack[0].File...)
			buf = append(buf, ':')
			buf = strconv.AppendInt(buf, int64(entry.Stack[0].Line), 10)
			buf = append(buf, ' ')
		}
	}
	// уровень записи
	level, ok := f.Levels[entry.Level]
	if !ok {
		level = entry.Level.String()
	}
	if level != "" {
		buf = append(buf, level...)
		buf = append(buf, ' ')
	}
	// категория
	if entry.Category != "" {
		buf = append(buf, '[')
		buf = append(buf, entry.Category...)
		buf = append(buf, "]: "...)
	}
	// основной текст
	if entry.Message != "" {
		buf = append(buf, entry.Message...)
	}
	// дополнительные поля
	for _, field := range entry.Fields {
		buf = append(buf, ' ')
		buf = append(buf, field.Name...)
		buf = append(buf, '=')
		switch value := field.Value.(type) {
		case string:
			buf = strconv.AppendQuote(buf, value)
		case []byte:
			buf = strconv.AppendQuoteToGraphic(buf, string(value))
		case error:
			buf = strconv.AppendQuote(buf, value.Error())
		case fmt.Stringer:
			buf = strconv.AppendQuote(buf, value.String())
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
		default:
			buf = append(buf, fmt.Sprint(value)...)
		}
	}
	buf = append(buf, '\n')
	return buf
}
