package log

import (
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Entry описывает запись в лог.
type Entry struct {
	Timestamp time.Time // временная метка
	Level     Level     // уровень
	Category  string    // название раздела
	Message   string    // текст
	Fields    []Field   // дополнительные поля
	Stack     []*Source // стек вызовов
	calldepth int       // уровень вложенности до исходного вызова
}

// NewEntry создает новое описание записи в лог.
func NewEntry(lvl Level, calldepth int, category, msg string, fields []Field) *Entry {
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
	entry.Stack = nil               // по умолчанию информация не заполняется
	entry.calldepth = calldepth + 2 // с учетом вызова этой функции и Encode
	return entry
}

// Free помещает объект для формирования записи лога обратно в пул.
func (e *Entry) Free() {
	entries.Put(e)
}

var emptySource = []*Source{}

// CallStack автоматически заполняет информацией о стеке вызовов, если она не
// была заполнена ранее.
func (e *Entry) CallStack(calldepth int) {
	if e.Stack != nil {
		return // уже заполнено
	}
	pc := make([]uintptr, 10)
	n := runtime.Callers(2+calldepth+e.calldepth, pc)
	if n == 0 {
		e.Stack = emptySource
		return // пустой стек, чтобы не заполнять еще раз
	}
	frames := runtime.CallersFrames(pc[:n])
	// TODO: pool []*Source
next:
	frame, more := frames.Next()
	if strings.HasPrefix(frame.Function, "runtime.") {
		return // не заполняем системными функциями
	}
	var source = &Source{
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

// Source возвращает строку с информацией об исходном файле.
func (e *Entry) Source(calldepth int) *Source {
	if e.Stack == nil {
		e.CallStack(calldepth + 1)
	}
	if len(e.Stack) == 0 {
		return nil
	}
	return e.Stack[0]
}

// Source описывает информацию об исходном файле с кодом.
type Source struct {
	Pkg  string // библиотека
	Func string // название функции
	File string // имя файла
	Line int    // номер строки
}

func (s *Source) String() string {
	return s.Pkg + "/" + s.File + ":" + strconv.Itoa(s.Line)
}
