package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Fields описывает список именованных полей события.
type Fields map[string]interface{}

// Entry описывает сообщение для записи в лог.
type Entry struct {
	Level   Level     `json:"level"`            // уровень
	Message string    `json:"msg,omitempty"`    // текст сообщения
	Fields  Fields    `json:"fields,omitempty"` // дополнительные поля
	start   time.Time // время старта для отложенных событий
	logger  *logger   // ссылка на лог
}

// set задает для сообщения уровень и текст. Если текст содержит перенос
// строки в конце, то он удаляется, чтобы не портить вывод сообщения с
// заполненными дополнительными полями.
func (e *Entry) set(level Level, msg string) *Entry {
	e.Level = level
	e.Message = strings.TrimSuffix(msg, "\n")
	return e
}

// log устанавливает для сообщения его уровень и текст, и отправляет его в лог.
func (e *Entry) handle() []error {
	return e.logger.handle(e)
}

// WithField добавляет к сообщению именованное поле с его значением.
func (e *Entry) WithField(name string, value interface{}) *Entry {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	e.Fields[name] = value
	return e
}

// WithFields добавляет к сообщению несколько именованных полей с их значениями.
func (e *Entry) WithFields(fields Fields) *Entry {
	if e.Fields == nil {
		e.Fields = make(map[string]interface{})
	}
	for name, value := range fields {
		e.Fields[name] = value
	}
	return e
}

// WithError добавляет к полям поле с сообщением об ошибке.
func (e *Entry) WithError(err error) *Entry {
	if err == nil {
		return e
	}
	return e.WithField("error", err.Error())
}

// WithSource добавляет в параметры сообщения название файла с исходным кодом
// и номер строки. calldepth указывает глубину анализа вызовов стека, а параметр
// fullpath - использовать полный путь к файлу.
func (e *Entry) WithSource(calldepth int, fullpath bool) *Entry {
	if _, file, line, ok := runtime.Caller(1 + calldepth); ok {
		if !fullpath {
			file = filepath.Base(file)
		}
		e.WithField("source", fmt.Sprintf("%s:%d", file, line))
	}
	return e
}

// Debug выводит в лог данное сообщение как отладочное.
func (e *Entry) Debug(v ...interface{}) {
	e.set(Debug, fmt.Sprintln(v...)).handle()
}

// Debugf выводит в лог данное сообщение как отладочное. Для формирования текста
// сообщения используется форматирование с помощью указанного шаблона.
func (e *Entry) Debugf(format string, v ...interface{}) {
	e.set(Debug, fmt.Sprintf(format, v...)).handle()
}

// Info выводит в лог данное сообщение как информационное.
func (e *Entry) Info(v ...interface{}) {
	e.set(Info, fmt.Sprintln(v...)).handle()
}

// Infof выводит в лог данное сообщение как информационное. Для формирования текста
// сообщения используется форматирование с помощью указанного шаблона.
func (e *Entry) Infof(format string, v ...interface{}) {
	e.set(Info, fmt.Sprintf(format, v...)).handle()
}

// Error выводит в лог данное сообщение как ошибочное.
func (e *Entry) Error(v ...interface{}) {
	e.set(Error, fmt.Sprintln(v...)).handle()
}

// Errorf выводит в лог данное сообщение как ошибочное. Для формирования текста
// сообщения используется форматирование с помощью указанного шаблона.
func (e *Entry) Errorf(format string, v ...interface{}) {
	e.set(Error, fmt.Sprintf(format, v...)).handle()
}

// copy возвращает независимую копию сообщения с установленным полем времени.
func (e *Entry) trace() *Entry {
	entry := *e              // создаем копию события
	entry.start = time.Now() // запоминаем время начала
	e.handle()               // отправляем уведомление
	return &entry            // возвращаем новое событие
}

// Trace выводит сообщение как информационное и возвращает его копию с
// установленным временем отсылки. Для данного типа сообщения можно вызвать
// метод Stop.
func (e *Entry) Trace(v ...interface{}) *Entry {
	return e.set(Info, fmt.Sprintln(v...)).trace()
}

// Trace формирует текст сообщения на основании указанного шаблона и выводит
// его как информационное. Возвращает копию сообщения с установленным временем
// отсылки, которое можно использовать с методом Stop.
func (e *Entry) Tracef(format string, v ...interface{}) *Entry {
	return e.set(Info, fmt.Sprintf(format, v...)).trace()
}

// Stop, в зависимости от ошибки, выводит либо информационное сообщение, либо
// сообщение об ошибке. Если это сообщение было сформированно с помощью метода
// Trace, то к полям сообщения автоматически добавляется поле duration.
func (e *Entry) Stop(err *error) {
	// устанавливаем уровень сообщения, в зависимости от ошибки
	if err == nil || *err == nil {
		e.Level = Info
	} else {
		e.Level = Error
		e.WithError(*err) // добавляем описание ошибки
	}
	// добавляем продолжительность выполнения
	if !e.start.IsZero() {
		e.WithField("duration", time.Since(e.start))
	}
	e.handle() // отправляем описание в лог
}

// entries содержит пул сообщений
var entries = sync.Pool{New: func() interface{} { return new(Entry) }}
