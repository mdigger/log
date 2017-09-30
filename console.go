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
	Levels     map[Level]string // переопределение строк для вывода уровня
}

// Encode форматирует в буфер запись лога для текстового консольного
// представления.
func (f Console) Encode(entry *Entry) []byte {
	var buf = buffer(buffers.Get().([]byte)[:0]) // получаем и сбрасываем буфер
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
		buf.WriteByte(' ')
	}
	// уровень записи
	level, ok := f.Levels[entry.Level]
	if !ok {
		level = entry.Level.String()
	}
	if level != "" {
		buf.WriteString(level)
		buf.WriteByte(' ')
	}
	// категория
	if entry.Category != "" {
		buf.WriteByte('[')
		buf.WriteString(entry.Category)
		buf.WriteString("]: ")
	}
	// основной текст
	if entry.Message != "" {
		buf.WriteString(entry.Message)
	}
	// дополнительные поля
	for _, field := range entry.Fields {
		buf.WriteByte(' ')
		buf.WriteString(field.Name)
		buf.WriteByte('=')
		switch value := field.Value.(type) {
		case nil:
			buf.WriteString("nil")
		case string:
			buf.WriteQuote(value)
		case []byte:
			buf = strconv.AppendQuoteToGraphic(buf, string(value))
		case error:
			buf.WriteQuote(value.Error())
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
			buf.WriteByte('"')
			if !value.IsZero() {
				buf = value.AppendFormat(buf, "2006-01-02 15:04:05")
			}
			buf.WriteByte('"')
		case fmt.Stringer:
			buf.WriteQuote(value.String())
		default:
			buf.WriteString(fmt.Sprint(value))
		}
	}
	buf.WriteByte('\n')
	return buf
}
