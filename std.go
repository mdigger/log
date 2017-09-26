package log

import "log"

// std обеспечивает поддержку стандартного лога.
type std struct {
	l   Logger // обработчик вывода в лог
	lvl Level  // уровень сообщений
}

// newStd возвращает стандартный лог, в качестве обработчика которого привязан
// наш лог. Все записи, которые будут сделаны в этот лог, будут автоматически
// перезаписаны в лог, обрабатываемый Handler. При этом будет использоваться
// заданный уровень и заданное имя категории.
func newStd(l Logger, lvl Level) *log.Logger {
	var std = &std{l: l, lvl: lvl}
	return log.New(std, "", 0)
}

// Write поддержка интерфейса io.Writer.
func (w *std) Write(p []byte) (int, error) {
	// убираем символ перехода на новую строку
	if l := len(p); l > 0 && p[l-1] == '\n' {
		p = p[:l-1]
	}
	// выводим запись в лог
	err := w.l.h.Write(w.lvl, 3, w.l.name, string(p), w.l.fields)
	return len(p), err
}
