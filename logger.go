package log

import (
	"log"
	"strconv"
)

// Field описывает именованное поле с дополнительными данными записи в лог.
type Field struct {
	Name  string
	Value interface{}
}

// Handler описывает интерфейс обработчика логов. В качестве параметров
// передаются уровень записи, глубина вложенных вызовов для определения
// исходного файла и строки с кодом, название категории, сообщения и список
// дополнительных полей.
type Handler interface {
	Write(lvl Level, calldepth int, name, msg string, fields []Field) error
}

// Logger описывает именованный раздел лога.
type Logger struct {
	name   string  // имя раздела
	h      Handler // обработчик лога
	fields []Field // дополнительные поля
}

// NewLogger возвращает новый Logger для записи лога с помощью обработчика h.
func NewLogger(h Handler) *Logger {
	return &Logger{h: h}
}

// New возвращает новую категорию с новым именем для того же лога. Имена
// добавляются через точку, используя предыдущее имя лога.
func (l *Logger) New(name string, fields ...interface{}) *Logger {
	if l.name != "" {
		name = l.name + "." + name
	}
	return &Logger{h: l.h, name: name, fields: l.appendFields(fields)}
}

// Log записывает в лог сообщение с заданным уровнем.
func (l *Logger) Log(lvl Level, msg string, fields ...interface{}) {
	l.h.Write(lvl, 1, l.name, msg, l.appendFields(fields))
}

// Trace записывает в лог низкоуровневое отладочное сообщение.
func (l *Logger) Trace(msg string, fields ...interface{}) {
	l.h.Write(TRACE, 1, l.name, msg, l.appendFields(fields))
}

// Debug записывает в лог низкоуровневое отладочное сообщение.
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.h.Write(DEBUG, 1, l.name, msg, l.appendFields(fields))
}

// Info записывает в лог информационное сообщение.
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.h.Write(INFO, 1, l.name, msg, l.appendFields(fields))
}

// Warn записывает в лог предупреждающее сообщение.
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.h.Write(WARN, 1, l.name, msg, l.appendFields(fields))
}

// Error записывает в лог сообщение об ошибке.
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.h.Write(ERROR, 1, l.name, msg, l.appendFields(fields))
}

// Fatal записывает в лог сообщение о критической ошибке.
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.h.Write(FATAL, 1, l.name, msg, l.appendFields(fields))
}

// With возвращает новую запись в лог с дополнительными параметрами.
func (l *Logger) With(fields ...interface{}) *Logger {
	return &Logger{h: l.h, name: l.name, fields: l.appendFields(fields)}
}

// StdLog возвращает стандартный лог.
func (l *Logger) StdLog(lvl Level) *log.Logger {
	return newStd(l, lvl)
}

// appendFields возвращает новый список дополнительных атрибутов, добавляя
// к уже существующим новые.
func (l *Logger) appendFields(fields []interface{}) []Field {
	switch len(fields) {
	case 0, 1:
		return l.fields
	case 2:
		name, ok := fields[0].(string)
		if !ok || name == "" {
			name = "key" + strconv.Itoa(len(l.fields)+1)
		}
		return append(l.fields, Field{Name: name, Value: fields[1]})
	}
	var list = make([]Field, len(fields)>>1)
	var name string
	for i, field := range fields {
		var pos = i >> 1
		if i%2 == 1 {
			list[pos] = Field{Name: name, Value: field}
			continue
		}
		var ok bool
		name, ok = field.(string)
		if !ok || name == "" {
			name = "key" + strconv.Itoa(len(l.fields)+pos+1)
		}
	}
	if l.fields == nil {
		return list
	}
	return append(l.fields, list...)
}
