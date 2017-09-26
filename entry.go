package log

import (
	"runtime"
	"strings"
	"time"
)

// Source описывает информацию об исходном файле с кодом.
type Source struct {
	Pkg  string // библиотека
	Func string // название функции
	File string // имя файла
	Line int    // номер строки
}

// Entry используется при форматировании записи лога.
type Entry struct {
	Timestamp time.Time // время инициализируется при записи в лог
	Level     Level     // уровень записи
	Category  string    // название лога
	Message   string    // сообщение
	Stack     []Source  // стек вызовов
	Fields    []Field   // дополнительные именованные поля
	calldepth int       // уровень вложенности до исходного вызова
}

// CallStack автоматически заполняет информацией о стеке вызовов, если она не
// была заполнена ранее.
func (e *Entry) CallStack(calldepth int) {
	if e.Stack != nil {
		return // уже заполнено
	}
	pc := make([]uintptr, 10)
	n := runtime.Callers(2+calldepth+e.calldepth, pc)
	if n == 0 {
		e.Stack = []Source{}
		return // пустой стек, чтобы не заполнять еще раз
	}
	frames := runtime.CallersFrames(pc[:n])
next:
	frame, more := frames.Next()
	if strings.HasPrefix(frame.Function, "runtime.") {
		return // не заполняем системными функциями
	}
	var source = Source{
		Line: frame.Line,
		Func: frame.Function,
		File: frame.File,
	}
	if lastslash := strings.LastIndex(source.Func, "/"); lastslash >= 0 {
		source.Pkg = source.Func[:lastslash] + "/"
		source.Func = source.Func[lastslash+1:]
	}
	if period := strings.Index(source.Func, "."); period >= 0 {
		source.Pkg += source.Func[:period]
		source.Func = source.Func[period+1:]
	}
	source.Func = strings.Replace(source.Func, "·", ".", -1)
	if lastslash := strings.LastIndex(source.File, "/"); lastslash >= 0 {
		source.File = source.File[lastslash+1:]
	}
	e.Stack = append(e.Stack, source)
	if more {
		goto next
	}
}
