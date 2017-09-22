package log

import "log"

// stdLog обеспечивает поддержку стандартного лога.
type stdLog struct {
	h    Handler // обработчик вывода в лог
	lvl  Level   // уровень сообщений
	name string  // имя категории
}

// newStdLog возвращает стандартный лог, в качестве обработчика которого привязан
// наш лог. Все записи, которые будут сделаны в этот лог, будут автоматически
// перезаписаны в лог, обрабатываемый Handler. При этом будет использоваться
// заданный уровень и заданное имя категории.
func newStdLog(h Handler, lvl Level, name string) *log.Logger {
	var std = &stdLog{h: h, lvl: lvl, name: name}
	return log.New(std, "", 0)
}

// Write поддержка интерфейса io.Writer.
func (w *stdLog) Write(p []byte) (int, error) {
	// убираем символ перехода на новую строку
	if l := len(p); l > 0 && p[l-1] == '\n' {
		p = p[:l-1]
	}
	// выводим запись в лог
	var err = w.h.Log(w.lvl, w.name, string(p))
	return len(p), err
}
