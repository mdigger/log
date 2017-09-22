package log

import (
	"errors"
	"log"
)

// Logger описывает именованный раздел (category) для вывода в лог. По умолчанию
// название раздела пустое.
type Logger struct {
	h    Handler // обработчик вывода в лог
	name string  // имя категории
}

// NewLogger возвращает новый лог с заданным обработчиком.
func NewLogger(handler Handler) *Logger {
	return &Logger{h: handler}
}

// New возвращает новую именованную категорию для вывода в исходный лог.
func (l *Logger) New(name string) *Logger {
	return &Logger{h: l.h, name: name}
}

// Log выводит в лог сообщение с указанным уровнем.
func (l *Logger) Log(lvl Level, msg string, fields ...interface{}) {
	l.h.Log(lvl, l.name, msg, fields...)
}

// Trace выводит необязательное отладочное сообщение.
func (l *Logger) Trace(msg string, fields ...interface{}) {
	l.h.Log(TRACE, l.name, msg, fields...)
}

// Debug выводит отладочное сообщение.
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.h.Log(DEBUG, l.name, msg, fields...)
}

// Info выводит информационное сообщение.
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.h.Log(INFO, l.name, msg, fields...)
}

// Warn выводит сообщение с предупреждением. Возвращает сформированную на
// основании текста сообщения ошибку, чтобы можно было ее использовать
// для дальнейшей обработки.
func (l *Logger) Warn(msg string, fields ...interface{}) error {
	l.h.Log(WARN, l.name, msg, fields...)
	return errors.New(msg)
}

// Error выводит сообщение об ошибке. Возвращает сформированную на основании
// текста сообщения ошибку, чтобы можно было ее использовать для дальнейшей
// обработки.
func (l *Logger) Error(msg string, fields ...interface{}) error {
	l.h.Log(ERROR, l.name, msg, fields...)
	return errors.New(msg)
}

// IfErr выводит в лог сообщение об ошибке, если ошибка err не пустая.
// При этом err автоматически добавляется как одно из дополнительных
// свойств в fields с именем "err". Возвращает исходную ошибку err без каких
// либо изменений.
func (l *Logger) IfErr(err error, msg string, fields ...interface{}) error {
	if err != nil {
		l.h.Log(ERROR, l.name, msg, append(fields, "err", err)...)
	}
	return nil
}

// Fatal выводит в лок высокоприоритетное сообщение об ошибке.
func (l *Logger) Fatal(msg string, fields ...interface{}) error {
	l.h.Log(FATAL, l.name, msg, fields...)
	return errors.New(msg)
}

// StdLog позволяет подменить стандартный лог на данный. В качестве параметров
// задается уровень формируемых записей в лог и имя раздела лога.
//
// Сделал эту функцию специально, потому что в некоторых случаях оказывается,
// что стандартный лог golang просто предопределен и переопределить его на
// что-то другое невозможно. Например, в http.Server. С помощью данного
// "костыля" это становится возможным. Но вызов методов SetFlag и добавление
// даты и времени может привести к тому, что время будет выводиться два раза
// и в разных форматах. А SetOutput вообще может порушить идиллию. Но тут уж
// ничего не поделаешь - стандартный Logger не является интерфейсом.
func (l *Logger) StdLog(lvl Level, name string) *log.Logger {
	return newStdLog(l.h, lvl, name)
}

// WithFields возвращает частичную запись в лог с заполненными полями.
func (l *Logger) WithFields(fields Fields) *Entry {
	return &Entry{logger: l, fields: fields}
}

// WithField возвращает частичную запись в лог с заполненным именованным полем.
func (l *Logger) WithField(name string, value interface{}) *Entry {
	return &Entry{logger: l, fields: Fields{name: value}}
}

// WithError возвращает частичную запись в лог с заполненным именованным полем
// с ошибкой.
func (l *Logger) WithError(err error) *Entry {
	var fields Fields
	if err != nil {
		fields = Fields{"err": err}
	}
	return &Entry{logger: l, fields: fields}
}
