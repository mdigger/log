package log

import (
	"runtime"
	"strconv"
	"strings"
)

// StackError описывает стандартную ошибку с добавлением информации о стеке
// вызовов.
type StackError struct {
	Err   error    // оригинальная ошибка
	Stack []Source // стек вызовов
}

// NewError формирует новую ошибку, добавляя информацию о стеке вызовов.
func NewError(err error) *StackError {
	var result = &StackError{Err: err}
	var pc [32]uintptr
	n := runtime.Callers(2, pc[:])
	if n == 0 {
		return result
	}
	result.Stack = make([]Source, 0, n-1)
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		if strings.HasPrefix(frame.Function, "runtime.") {
			break // не заполняем системными функциями
		}
		result.Stack = append(result.Stack, newSource(frame))
		if !more {
			break
		}
	}
	return result
}

// Error возвращает строковое описание ошибки.
func (e *StackError) Error() string {
	return e.Err.Error()
}

// Source описывает информацию об исходном файле с кодом.
type Source struct {
	Pkg  string // библиотека
	Func string // название функции
	File string // имя файла
	Line int    // номер строки
}

// NewSource возвращает описание об исходном файле и функции.
func newSource(frame runtime.Frame) Source {
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
	return source
}

// String возвращает строку с исходным файлом и номером строки.
func (s Source) String() string {
	return s.Pkg + "/" + s.File + ":" + strconv.Itoa(s.Line)
}
