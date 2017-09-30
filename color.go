package log

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mdigger/errors"
)

// Color выводит лог в более удобном для чтения в консоли виде, используя
// цветовые выделения и помещая параметры на новую строку.
type Color struct {
	Levels    map[Level]string // переопределение строк для вывода уровня
	KeyIndent int              // отступ от значения дополнительного параметра
}

// Encode форматирует в буфер запись лога для текстового консольного
// представления.
func (f Color) Encode(entry *Entry) []byte {
	var buf = buffer(buffers.Get().([]byte)[:0]) // получаем и сбрасываем буфер
	// выводим время
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	buf.WriteString("\x1b[2m")
	buf = entry.Timestamp.AppendFormat(buf, "15:04:05.000000")
	buf.WriteString("\x1b[0m ")
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
		buf.WriteString("\x1b[7m\x1b[")
		switch entry.Level & -32 {
		case FATAL:
			buf.WriteString("35")
		case ERROR:
			buf.WriteString("91")
		case WARN:
			buf.WriteString("93")
		case INFO:
			buf.WriteString("92")
		case DEBUG:
			buf.WriteString("94")
		case TRACE:
			buf.WriteString("96")
		default:
			buf.WriteString("37")
		}
		buf.WriteByte('m')
		buf.WriteString(level)
		buf.WriteString("\x1b[0m ")
	}
	// категория
	if entry.Category != "" {
		buf.WriteString("\x1b[2m[\x1b[0m\x1b[92m")
		buf.WriteString(entry.Category)
		buf.WriteString("\x1b[0m\x1b[2m]:\x1b[0m ")
	}
	// основной текст
	if entry.Message != "" {
		buf.WriteString(entry.Message)
	}
	// дополнительные поля
	for _, field := range entry.Fields {
		buf.WriteString("\n    \x1b[36m")
		buf.WriteString(field.Name)
		buf.WriteString("\x1b[0m")
		for i := 0; i < f.KeyIndent-len(field.Name); i++ {
			buf.WriteByte(' ')
		}
		buf.WriteString("\x1b[2m=\x1b[0m")
		if f.KeyIndent > 0 {
			buf.WriteByte(' ')
		}
		switch value := field.Value.(type) {
		case nil:
			buf.WriteString("nil")
		case string:
			buf.WriteString(value)
		case error:
			buf.WriteQuote(value.Error())
			if value, ok := value.(*errors.Error); ok {
				if cause := value.Cause(); cause != nil {
					buf.WriteString(" \x1b[2mcause: \x1b[0m\x1b[91m")
					fmt.Fprintf(&buf, "%#v", cause)
					buf.WriteString("\x1b[0m")
				}
				for _, src := range value.Stacks() {
					buf.WriteString("\n\t- ")
					buf.WriteString(src.Func)
					buf.WriteString(" \x1b[2m[")
					buf.WriteString(src.String())
					buf.WriteString("]\x1b[0m")
				}
			} else {
				buf.WriteString(" \x1b[91m")
				fmt.Fprintf(&buf, "%#v", value)
				buf.WriteString("\x1b[0m")
			}
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
			buf.WriteString(value.String())
		default:
			buf.WriteString(fmt.Sprint(value))
		}
	}
	// // для ошибок выводим стек вызовов
	// if entry.Level >= WARN {
	// 	if entry.Stack == nil {
	// 		entry.CallStack(1)
	// 	}
	// 	for _, src := range entry.Stack {
	// 		buf = append(buf, "\n  \x1b[2m- "...)
	// 		buf = append(buf, src.Pkg...)
	// 		buf = append(buf, "/\x1b[0m"...)
	// 		buf = append(buf, src.File...)
	// 		buf = append(buf, "\x1b[2m:\x1b[0m"...)
	// 		buf = strconv.AppendInt(buf, int64(src.Line), 10)
	// 		buf = append(buf, " \x1b[2m(\x1b[0m\x1b[36m"...)
	// 		buf = append(buf, src.Func...)
	// 		buf = append(buf, "\x1b[0m\x1b[2m)\x1b[0m"...)
	// 	}
	// }
	buf.WriteByte('\n')
	return buf
}
