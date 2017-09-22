package log

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// JSON отвечает за форматирование лога в формате JSON. В качестве значения
// может быть задана строка, которая будет использоваться для отступов
// при форматировании.
type JSON string

var _ StreamFormatter = new(JSON) // проверяем поддержу интерфейса

// Log осуществляет форматирование записи лога в формат JSON.
func (f JSON) Log(w io.Writer, lvl Level, category, msg string,
	fields ...interface{}) error {
	var enc = json.NewEncoder(w)
	if f != "" {
		enc.SetIndent("", string(f))
	}
	var entry = entries.Get().(*jsonEntry)
	entry.Timestamp = time.Now().Unix()
	entry.Level = lvl
	entry.Category = category
	entry.Message = msg
	switch len(fields) {
	case 0: // нет дополнительных полей
		break
	case 1: // дополнительные поля представлены одним элементом
		if list, ok := fields[0].(map[string]interface{}); ok {
			entry.Fields = make(map[string]json.RawMessage, len(list))
			for name, value := range list {
				entry.Fields[name] = fieldValue(value)
			}
		}
	default:
		entry.Fields = make(map[string]json.RawMessage, len(fields)<<1)
		var name string
		for i, field := range fields {
			if i%2 == 0 {
				if s, ok := field.(string); ok {
					name = s
				} else {
					name = fmt.Sprint(name)
				}
			} else {
				entry.Fields[name] = fieldValue(field)
			}
		}
	}
	var err = enc.Encode(entry)
	entries.Put(entry)
	return err
}

// jsonEntry описывает структуру записи лога для записи в формате JSON.
type jsonEntry struct {
	Timestamp int64                      `json:"time"`
	Level     Level                      `json:"lvl"`
	Category  string                     `json:"category,omitempty"`
	Message   string                     `json:"msg"`
	Fields    map[string]json.RawMessage `json:"keys,omitempty"`
}

var entries = sync.Pool{New: func() interface{} { return new(jsonEntry) }}

// fieldValue трансформирует некоторые значения в вид, удобный для JSON.
func fieldValue(value interface{}) json.RawMessage {
	if err, ok := value.(error); ok {
		value = fmt.Sprintf("[%T]: %[1]s", err)
	}
repeat:
	data, err := json.Marshal(value)
	if err != nil {
		value = fmt.Sprint(value)
		goto repeat
	}
	return data
}
