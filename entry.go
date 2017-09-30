package log

import (
	"sync"
	"time"
)

// Entry описывает запись в лог.
type Entry struct {
	Timestamp time.Time // временная метка
	Level     Level     // уровень
	Category  string    // название раздела
	Message   string    // текст
	Fields    []Field   // дополнительные поля
}

// NewEntry создает новое описание записи в лог.
func NewEntry(lvl Level, category, msg string, fields []Field) *Entry {
	var names = make(map[string]int, len(fields))
	var result = make([]Field, 0, len(fields))
	for i, field := range fields {
		if field.Name == "" {
			field.Name = "_" // подменяем пустое имя
		}
		// проверяем, что поле с таким именем уже было
		if pos, ok := names[field.Name]; ok {
			result[pos].Value = field.Value // заменяем старое значение на новое
			continue
		}
		result = append(result, field)
		names[field.Name] = i // сохраняем позицию
	}
	var entry = entries.Get().(*Entry)
	entry.Timestamp = time.Time{} // не устанавливаем время до записи
	entry.Level = lvl
	entry.Category = category
	entry.Message = msg
	entry.Fields = result
	return entry
}

// Free помещает объект для формирования записи лога обратно в пул.
func (e *Entry) Free() {
	entries.Put(e)
}

var entries = sync.Pool{New: func() interface{} { return new(Entry) }}
