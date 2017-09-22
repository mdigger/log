package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// SourceInfo описывает информацию о файле с исходным кодом, в котором произошел
// вызов.
type SourceInfo struct {
	Pkg  string `json:"pkg"`  // название библиотеки
	File string `json:"file"` // имя файла
	Line int    `json:"line"` // номер строки в файле
	Func string `json:"func"` // название функции
}

// String возвращает строковое представление информации об исходном файле и
// номере строки.
func (si *SourceInfo) String() string {
	return fmt.Sprintf("%s/%s:%02d (%s)", si.Pkg, si.File, si.Line, si.Func)
}

// Source возвращает имя текущего файла с исходным кодом, номер строки, название
// пакета и функции.
//
// Данную функцию удобно использовать для добавления информации об исходном
// файле в лог:
// 	log.Warning("warning", "src", log.Source(0))
func Source(skip int) *SourceInfo {
	var pc, file, line, ok = runtime.Caller(skip + 1)
	if !ok {
		return nil
	}
	file = filepath.Base(file)
	var name = runtime.FuncForPC(pc).Name()
	var pkg string
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}
	name = strings.Replace(name, "·", ".", -1)
	return &SourceInfo{
		File: file,
		Line: line,
		Pkg:  pkg,
		Func: name,
	}
}

// CallStack возвращает информацию о текущем стеке вызовов функций. Может
// использоваться для добавления информации о стеке вызовов в дополнительные
// параметры лога:
// 	log.Debug("call stack", "stack", CallStack(0))
func CallStack(skip int) []*SourceInfo {
	var pc = make([]uintptr, 10)
	var n = runtime.Callers(2+skip, pc)
	if n == 0 {
		return nil
	}
	pc = pc[:n]
	var frames = runtime.CallersFrames(pc)
	var result = make([]*SourceInfo, 0, n)
	for {
		frame, more := frames.Next()
		var name = frame.Function
		if strings.HasPrefix(name, "runtime.") {
			break
		}
		var pkg string
		if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
			pkg += name[:lastslash] + "/"
			name = name[lastslash+1:]
		}
		if period := strings.Index(name, "."); period >= 0 {
			pkg += name[:period]
			name = name[period+1:]
		}
		name = strings.Replace(name, "·", ".", -1)
		result = append(result, &SourceInfo{
			File: filepath.Base(frame.File),
			Line: frame.Line,
			Pkg:  pkg,
			Func: name,
		})
		if !more {
			break
		}
	}
	return result
}
