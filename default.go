package log

import (
	"errors"
	"io"
	"log"
	"os"
)

// defaultHandler является предопределенным консольным обработчиком лога.
var defaultHandler = NewStreamHandler(os.Stderr, DEBUG,
	&Console{TimeFormat: "2006-01-02 15:04:05"})

// SetLevel изменяет уровень фильтра для вывода сообщений в лог по умолчанию.
// Изначально фильтр установлен в DEBUG.
func SetLevel(lvl Level) {
	defaultHandler.SetLevel(lvl)
}

// SetOutput переопределяет вывод лога по умолчанию. Изначально используется
// os.Stderr.
func SetOutput(w io.Writer) {
	defaultHandler.SetOutput(w)
}

// SetFormat переопределяет настройки формата вывода лога по умолчанию.
// Изначально из настроек задан только формат временной метки
// "2006-01-02 15:04:05". Все остальные настройки оставлены по умолчанию.
func SetFormat(format StreamFormatter) {
	defaultHandler.SetFormat(format)
}

// IsTTY возвращает true, если лог выводится в терминал или в файл.
func IsTTY() bool {
	return defaultHandler.IsTTY()
}

// New возвращает новый именованный раздел лога по умолчанию.
func New(name string) *Logger {
	return &Logger{h: defaultHandler, name: name}
}

// Log выводит сообщение с указанным уровнем в лог по умолчанию.
func Log(lvl Level, msg string, fields ...interface{}) {
	defaultHandler.Log(lvl, "", msg, fields...)
}

// Trace выводит необязательное отладочное сообщение в лог по умолчанию.
func Trace(msg string, fields ...interface{}) {
	defaultHandler.Log(TRACE, "", msg, fields...)
}

// Debug выводит отладочное сообщение в лог по умолчанию.
func Debug(msg string, fields ...interface{}) {
	defaultHandler.Log(DEBUG, "", msg, fields...)
}

// Info выводит информационное сообщение в лог по умолчанию.
func Info(msg string, fields ...interface{}) {
	defaultHandler.Log(INFO, "", msg, fields...)
}

// Warn выводит сообщение с предупреждением в лог по умолчанию. Возвращает
// сформированную на основании текста сообщения ошибку, чтобы можно было ее
// использовать для дальнейшей обработки.
func Warn(msg string, fields ...interface{}) error {
	defaultHandler.Log(WARN, "", msg, fields...)
	return errors.New(msg)
}

// Error выводит сообщение об ошибке в лог по умолчанию. Возвращает
// сформированную на основании текста сообщения ошибку, чтобы можно было ее
// использовать для дальнейшей обработки.
func Error(msg string, fields ...interface{}) error {
	defaultHandler.Log(ERROR, "", msg, fields...)
	return errors.New(msg)
}

// IfErr выводит  в лог по умолчанию сообщение об ошибке, если ошибка err не
// пустая. При этом err автоматически добавляется как одно из дополнительных
// свойств в fields с именем "err". Возвращает исходную ошибку err без каких
// либо изменений.
func IfErr(err error, msg string, fields ...interface{}) error {
	if err != nil {
		defaultHandler.Log(ERROR, "", msg, append(fields, "err", err)...)
	}
	return nil
}

// Fatal выводит сообщение о приоритетной ошибке в лог по умолчанию. Возвращает
// сформированную на основании текста сообщения ошибку, чтобы можно было ее
// использовать для дальнейшей обработки.
func Fatal(msg string, fields ...interface{}) error {
	defaultHandler.Log(FATAL, "", msg, fields...)
	return errors.New(msg)
}

// StdLogger позволяет подменить стандартный лог на лог по умолчанию. В качестве
// параметров задается уровень формируемых записей в лог и имя раздела лога.
func StdLogger(lvl Level, name string) *log.Logger {
	return newStdLog(defaultHandler, lvl, name)
}

// defaultLogger является предопределенным консольным логом по умолчанию.
var defaultLogger = NewLogger(defaultHandler)

// WithFields возвращает частичную запись в лог с заполненными полями. Запись
// будет относиться к логу по умолчанию.
func WithFields(fields Fields) *Entry {
	return &Entry{logger: defaultLogger, fields: fields}
}

// WithField возвращает частичную запись в лог с заполненным именованным полем.
// Запись будет относиться к логу по умолчанию.
func WithField(name string, value interface{}) *Entry {
	return &Entry{logger: defaultLogger, fields: Fields{name: value}}
}

// WithError возвращает частичную запись в лог с заполненным именованным полем
// с ошибкой. Запись будет относиться к логу по умолчанию.
func WithError(err error) *Entry {
	var fields Fields
	if err != nil {
		fields = Fields{"err": err}
	}
	return &Entry{logger: defaultLogger, fields: fields}
}
