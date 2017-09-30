package log

// std обеспечивает поддержку стандартного лога.
type std struct {
	l   *Logger // обработчик вывода в лог
	lvl Level   // уровень сообщений
}

// Write поддержка интерфейса io.Writer.
func (w *std) Write(p []byte) (l int, err error) {
	l = len(p)
	if l > 0 && p[l-1] == '\n' {
		p = p[:l-1] // убираем символ перехода на новую строку
	}
	err = w.l.h.Write(w.lvl, w.l.name, string(p), w.l.fields)
	return
}
