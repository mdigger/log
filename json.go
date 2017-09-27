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
func (f JSON) Encode(buf []byte, entry *Entry) []byte {
	buf = append(buf, `{"ts":`...)
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	buf = strconv.AppendInt(buf, entry.Timestamp.Unix(), 10)
	buf = append(buf, `,"lvl":`...)
	buf = strconv.AppendInt(buf, int64(entry.Level), 10)
	if entry.Category != "" {
		buf = append(buf, `,"log":`...)
		buf = strconv.AppendQuote(buf, entry.Category)
	}
	if entry.Message != "" {
		buf = append(buf, `,"msg":`...)
		buf = strconv.AppendQuote(buf, entry.Message)
	}
	for _, field := range entry.Fields {
		buf = append(buf, ',')
		buf = strconv.AppendQuote(buf, field.Name)
		buf = append(buf, ':')
		switch value := field.Value.(type) {
		case nil:
			buf = append(buf, "null"...)
		case string:
			buf = strconv.AppendQuote(buf, value)
		case []byte:
			b64 := make([]byte, base64.StdEncoding.EncodedLen(len(value)))
			base64.StdEncoding.Encode(b64, value)
			buf = append(buf, '"')
			buf = append(buf, b64...)
			buf = append(buf, '"')
		case error:
			if value == nil {
				buf = append(buf, "null"...)
			} else {
				buf = strconv.AppendQuote(buf, value.Error())
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
				buf = append(buf, `""`...)
			} else {
				buf = value.AppendFormat(buf, time.RFC3339)
			}
		case time.Duration:
			buf = strconv.AppendInt(buf, int64(value), 10)
		case fmt.Stringer:
			buf = strconv.AppendQuote(buf, value.String())
		default:
			if data, err := json.Marshal(value); err == nil {
				buf = append(buf, data...)
			} else {
				buf = strconv.AppendQuote(buf, fmt.Sprint(value))
			}
		}
	}
	// для предупреждений и ошибок добавляем информацию об исходном файле
	if entry.Level >= WARN {
		if src := entry.Source(1); src != nil {
			buf = append(buf, `,"@src":`...)
			buf = append(buf, '"')
			buf = append(buf, src.Pkg...)
			buf = append(buf, '/')
			buf = append(buf, src.File...)
			buf = append(buf, ':')
			buf = strconv.AppendInt(buf, int64(src.Line), 10)
			buf = append(buf, '"')
		}
	}
	buf = append(buf, "}\n"...)
	return buf
}
