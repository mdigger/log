package log

import (
	"flag"
	"io"
	"log"
	"os"
)

// устанавливаем формат для лога по умолчанию из переменной окружения.
func init() {
	if config, ok := os.LookupEnv("LOG"); ok {
		h.Set(config)
	}
}

// default используется как лог по умолчанию.
var h = &Writer{w: os.Stderr, lvl: DEBUG, enc: &Console{
	TimeFormat: "2006-01-02 15:04:05",
}}

// Flag возвращает лог по умолчанию в качестве значения для установки через
// параметры приложения.
func Flag() flag.Value {
	return h
}

// SetLevel изменяет уровень фильтра для вывода сообщений в лог по умолчанию.
// Изначально фильтр установлен в DEBUG.
func SetLevel(lvl Level) {
	h.SetLevel(lvl)
}

// SetOutput переопределяет вывод лога по умолчанию. Изначально используется
// os.Stderr.
func SetOutput(w io.Writer) {
	h.SetOutput(w)
}

// SetFormat переопределяет настройки формата вывода лога по умолчанию.
// Изначально из настроек задан только формат временной метки
// "2006-01-02 15:04:05". Все остальные настройки оставлены по умолчанию.
func SetFormat(enc Encoder) {
	h.SetFormat(enc)
}

// IsTTY возвращает true, если лог выводится в терминал или в файл.
func IsTTY() bool {
	return h.IsTTY()
}

// New возвращает новый именованный раздел лога по умолчанию.
func New(name string, fields ...interface{}) *Logger {
	return &Logger{h: h, name: name, fields: h.with(fields)}
}

// Log выводит сообщение с указанным уровнем в лог по умолчанию.
func Log(lvl Level, msg string, fields ...interface{}) {
	h.Write(lvl, 1, "", msg, h.with(fields))
}

// Trace выводит необязательное отладочное сообщение в лог по умолчанию.
func Trace(msg string, fields ...interface{}) {
	h.Write(TRACE, 1, "", msg, h.with(fields))
}

// Debug выводит отладочное сообщение в лог по умолчанию.
func Debug(msg string, fields ...interface{}) {
	h.Write(DEBUG, 1, "", msg, h.with(fields))
}

// Info выводит информационное сообщение в лог по умолчанию.
func Info(msg string, fields ...interface{}) {
	h.Write(INFO, 1, "", msg, h.with(fields))
}

// Warn выводит сообщение с предупреждением в лог по умолчанию.
func Warn(msg string, fields ...interface{}) {
	h.Write(WARN, 1, "", msg, h.with(fields))
}

// Error выводит сообщение об ошибке в лог по умолчанию.
func Error(msg string, fields ...interface{}) {
	h.Write(ERROR, 1, "", msg, h.with(fields))
}

// Fatal выводит сообщение о критической ошибке в лог по умолчанию.
func Fatal(msg string, fields ...interface{}) {
	h.Write(FATAL, 1, "", msg, h.with(fields))
}

// With возвращает новую запись в лог с дополнительными параметрами.
func With(fields ...interface{}) *Logger {
	return &Logger{h: h, name: "", fields: h.with(fields)}
}

// StdLog возвращает обертку лога в стандартный.
func StdLog(lvl Level, name string, fields ...interface{}) *log.Logger {
	return log.New(&std{l: &Logger{h: h, name: name, fields: h.with(fields)},
		lvl: lvl}, "", 0)
}
