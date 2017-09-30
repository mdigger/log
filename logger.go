package log

import (
	"fmt"
	"log"
)

// Handler описывает интерфейс для записи лога.
type Handler interface {
	Write(lvl Level, category, msg string, fields []Field) error
}

// Field описывает дополнительное именованное поле лога.
type Field struct {
	Name  string
	Value interface{}
}

// Fields описывает список дополнительных полей.
type Fields = map[string]interface{}

// Logger описывает именованный раздел лога.
type Logger struct {
	h      Handler // обработчик лога
	name   string  // название раздела
	fields []Field // дополнительные поля
}

// NewLogger возвращает новый лог с указанным обработчиком.
func NewLogger(h Handler) *Logger {
	return &Logger{h: h}
}

// New возвращает новый именованный раздел лога. Новое имя будет добавлено к
// имени предыдущего раздела лога с разделителем ".".
func (l *Logger) New(name string, fields ...interface{}) *Logger {
	if name == "" {
		name = l.name
	} else if l.name != "" {
		name = l.name + "." + name
	}
	return &Logger{
		h:      l.h,
		name:   name,
		fields: l.with(fields),
	}
}

// Log добавляет запись в лог с указанным уровнем.
func (l *Logger) Log(lvl Level, msg string, fields ...interface{}) {
	l.h.Write(lvl, l.name, msg, l.with(fields))
}

// Trace записывает в лог сообщение с уровнем ниже отладочного.
func (l *Logger) Trace(msg string, fields ...interface{}) {
	l.h.Write(TRACE, l.name, msg, l.with(fields))
}

// Debug записывает в лог отладочное сообщение.
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.h.Write(DEBUG, l.name, msg, l.with(fields))
}

// Info записывает в лог информационное сообщение.
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.h.Write(INFO, l.name, msg, l.with(fields))
}

// Warn записывает в лог сообщение с предупреждением.
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.h.Write(WARN, l.name, msg, l.with(fields))
}

// Error записывает в лог сообщение с ошибкой.
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.h.Write(ERROR, l.name, msg, l.with(fields))
}

// Fatal записывает в лог сообщение с критической ошибкой.
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.h.Write(FATAL, l.name, msg, l.with(fields))
}

// StdLog возвращает обертку лога в стандартный. В качестве параметров
// указывается уровень сообщений, который будет использоваться по умолчанию
// для всех записей лога.
func (l *Logger) StdLog(lvl Level) *log.Logger {
	return log.New(&std{l: l, lvl: lvl}, "", 0)
}

// With возвращает запись лога с частично заполненными дополнительными полями.
// В качестве именованных параметров можно использовать последовательное
// перечисление имен (строк) и соответствующих значений или непосредственно
// список имен со значениями в виде Field. Отдельно обрабатываются ошибки,
// которые тоже могут быть переданы без имени параметра: в этом случае будет
// использовано имя "error", если ошибка не пустая. Если вы ошиблись и для
// последнего элемента не задали значение, то такой элемент будет
// проигнорирован. Эти правила действительны и для всех методов Logger.
func (l *Logger) With(fields ...interface{}) *Logger {
	return &Logger{
		h:      l.h,
		name:   l.name,
		fields: l.with(fields),
	}
}

// with при любом изменении полей возвращает их объединенную копию. В противном
// случае возвращает список как есть.
func (l *Logger) with(fields []interface{}) []Field {
	if len(fields) == 0 {
		return l.fields // ничего нового не будет
	}
	var result []Field
	// обрабатываем новые поля, добавляя их в список
	for i := 0; i < len(fields); i++ {
		var name string
		switch val := fields[i].(type) {
		case Field:
			result = append(result, val)
			continue
		case []Field:
			result = append(result, val...)
			continue
		case Fields: // поля уже является самостоятельным списком
			for name, value := range val {
				result = append(result, Field{name, value})
			}
			continue
		case error: // для ошибок без имени поля используем поле "error"
			if val != nil {
				result = append(result, Field{"error", val})
			}
			continue
		case string: // название поля
			name = val
		case fmt.Stringer:
			name = val.String()
		default: // не известный тип для названия поля
			name = fmt.Sprint(val)
		}
		if i == len(fields)-1 {
			break // это последний элемент в списке и значения для него нет
		}
		i++ // увеличиваем счетчик прочитанных
		// читаем следующее значение в списке
		result = append(result, Field{name, fields[i]})
	}
	return append(l.fields, result...)
}
