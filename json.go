package log

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// JSON формирует запись в лог в формате JSON.
type JSON struct{}

// Encode возвращает представление записи в лог в формате JSON.
func (f JSON) Encode(entry *Entry) []byte {
	var buf = buffer(buffers.Get().([]byte)[:0]) // получаем и сбрасываем буфер
	buf.WriteString(`{"ts":`)
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	buf = strconv.AppendInt(buf, entry.Timestamp.Unix(), 10)
	buf.WriteString(`,"lvl":`)
	buf = strconv.AppendInt(buf, int64(entry.Level), 10)
	if entry.Category != "" {
		buf.WriteString(`,"log":`)
		buf.WriteQuote(entry.Category)
	}
	if entry.Message != "" {
		buf.WriteString(`,"msg":`)
		buf.WriteQuote(entry.Message)
	}
	for _, field := range entry.Fields {
		buf.WriteByte(',')
		buf.WriteQuote(field.Name)
		buf.WriteByte(':')
		switch value := field.Value.(type) {
		case nil:
			buf.WriteString("null")
		case string:
			buf.WriteQuote(value)
		case []byte:
			buf.WriteQuote(base64.StdEncoding.EncodeToString(value))
		case error:
			if value == nil {
				buf.WriteString("null")
			} else {
				buf.WriteQuote(value.Error())
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
			if value.IsZero() {
				buf.WriteString(`""`)
			} else {
				buf = value.AppendFormat(buf, time.RFC3339)
			}
		case time.Duration:
			buf = strconv.AppendInt(buf, int64(value), 10)
		case fmt.Stringer:
			buf.WriteQuote(value.String())
		default:
			if data, err := json.Marshal(value); err == nil {
				buf = append(buf, data...)
			} else {
				buf.WriteQuote(fmt.Sprint(value))
			}
		}
	}
	// // для предупреждений и ошибок добавляем информацию об исходном файле
	// if entry.Level >= WARN {
	// 	if src := entry.Source(1); src != nil {
	// 		buf = append(buf, `,"@src":`...)
	// 		buf = append(buf, '"')
	// 		buf = append(buf, src.Pkg...)
	// 		buf = append(buf, '/')
	// 		buf = append(buf, src.File...)
	// 		buf = append(buf, ':')
	// 		buf = strconv.AppendInt(buf, int64(src.Line), 10)
	// 		buf = append(buf, '"')
	// 	}
	// }
	buf.WriteString("}\n")
	return buf
}
