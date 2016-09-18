package log

import (
	"fmt"
	"sync"
)

// Logger описывает интерфейс, который обеспечивает вывод в лог информации.
type Logger interface {
	SetLevel(Level)
	SetFlags(int)
	Println(level Level, v ...interface{})
	Printf(level Level, format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	WithField(name string, value interface{}) *Entry
	WithFields(fields Fields) *Entry
	WithError(err error) *Entry
	WithSource(calldepth int, short bool) *Entry
	Trace(v ...interface{}) *Entry
	Tracef(format string, v ...interface{}) *Entry
}

// Handler описывает интерфейс обработчика лога.
type Handler interface {
	Handle(*Entry) error
	SetFlags(int)
}

// New возвращает новый инициализированный лог с заданным обработчиком
// вывода.
func New(handlers ...Handler) Logger {
	return &logger{handlers: handlers}
}

// logger представляет внутреннюю реализацию логики работы лога.
type logger struct {
	Level
	handlers []Handler
	mu       sync.Mutex
}

// SetLevel задает минимальный уровень сообщения необходимый для вывода в лог.
func (l *logger) SetLevel(level Level) {
	l.mu.Lock()
	l.Level = level
	l.mu.Unlock()
}

// AddHandler добавляет в лог новый обработчик для вывода лога.
func (l *logger) AddHandler(handler Handler) {
	l.mu.Lock()
	if l.handlers == nil {
		l.handlers = []Handler{handler}
	} else {
		l.handlers = append(l.handlers, handler)
	}
	l.mu.Unlock()
}

// SetFlags устанавливает флаги для всех обработчиков лога.
func (l *logger) SetFlags(flag int) {
	l.mu.Lock()
	for _, handler := range l.handlers {
		handler.SetFlags(flag)
	}
	l.mu.Unlock()
}

// handle проверяет, что лог поддерживает сообщения такого уровня и, если да, то
// отправляет его в лог. После этого помещает сообщение в буфер для возможного
// дальнейшего использования.
func (l *logger) handle(e *Entry) (errors []error) {
	defer entries.Put(e) // освобождаем память по окончании
	// игнорируем, если уровень сообщения для лога ниже заданного
	if e.Level < l.Level {
		return nil
	}
	// вызываем обработчик вывода в лог
	l.mu.Lock()
	for _, handler := range l.handlers {
		if err := handler.Handle(e); err != nil {
			if errors == nil {
				errors = []error{err}
			} else {
				errors = append(errors, err)
			}
		}
	}
	l.mu.Unlock()
	return errors
}

// new создает новое сообщение для записи в лог для конкретного лога.
func (l *logger) new(level Level, msg string) *Entry {
	entry := entries.Get().(*Entry) // получаем из буфера
	// очищаем возможные заполненные поля, сохранившиеся в буфере
	entry.logger = l
	entry.Fields = nil
	return entry.set(level, msg)
}

// Println выводит сообщение с указанным уровнем в лог.
func (l *logger) Println(level Level, v ...interface{}) {
	l.new(level, fmt.Sprintln(v...)).handle()
}

// Printf формирует новое сообщение с указанным уровнем на основе указанного
// формата и выводит его в лог.
func (l *logger) Printf(level Level, format string, v ...interface{}) {
	l.new(level, fmt.Sprintf(format, v...)).handle()
}

// Debug выводит отладочное сообщение в лог.
func (l *logger) Debug(v ...interface{}) {
	l.new(Debug, fmt.Sprintln(v...)).handle()
}

// Debugf формирует новое сообщение для отладки на основе указанного формата и
// выводит его в лог.
func (l *logger) Debugf(format string, v ...interface{}) {
	l.new(Debug, fmt.Sprintf(format, v...)).handle()
}

// Info выводит стандартное сообщение в лог.
func (l *logger) Info(v ...interface{}) {
	l.new(Info, fmt.Sprintln(v...)).handle()
}

// Infof формирует новое стандартное сообщение на основе указанного формата и
// выводит его в лог.
func (l *logger) Infof(format string, v ...interface{}) {
	l.new(Info, fmt.Sprintf(format, v...)).handle()
}

// Error выводит в лог сообщение об ошибке.
func (l *logger) Error(v ...interface{}) {
	l.new(Error, fmt.Sprintln(v...)).handle()
}

// Errorf создает новое сообщение об ошибке на основании шаблона и выводит его
// в лог.
func (l *logger) Errorf(format string, v ...interface{}) {
	l.new(Error, fmt.Sprintf(format, v...)).handle()
}

func (l *logger) Trace(v ...interface{}) *Entry {
	return l.new(Info, fmt.Sprintln(v...)).trace()
}

func (l *logger) Tracef(format string, v ...interface{}) *Entry {
	return l.new(Info, fmt.Sprintf(format, v...)).trace()
}

// WithField возвращает инициализированную запись с указанным полем и значением
// для дальнейшей записи в лог.
func (l *logger) WithField(name string, value interface{}) *Entry {
	return l.new(Info, "").WithField(name, value)
}

// WithField возвращает инициализированную запись с указанными полями и
// значениями для дальнейшей записи в лог.
func (l *logger) WithFields(fields Fields) *Entry {
	return l.new(Info, "").WithFields(fields)
}

// WithError возвращает инициализированную запись с заполненным полем error для
// дальнейшей записи в лог.
func (l *logger) WithError(err error) *Entry {
	return l.new(Info, "").WithError(err)
}

// WithSource добавляет в параметры название файла с исходным кодом и номер
// строки. Параметр calldepth указывает глубину анализа стека вызовов, а
// fullpath — использовать полный путь к файлу.
func (l *logger) WithSource(calldepth int, fullpath bool) *Entry {
	return l.new(Info, "").WithSource(calldepth+1, fullpath)
}
