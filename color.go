package log

// Color определяет вывод для консоли в цвете. Обычно используется для
// вывода лога в консоль для отладки, т.к. является более читаемым, чем
// стандартный текст, за счет цветового выделения и того, что каждый
// дополнительны параметр записи выводится с новой строки с минимальным
// отступом.
var Color = &Console{
	TimeFormat: "15:04:05",
	UTC:        false,
	Time:       Pair{"\033[2m", "\033[22m "},
	Levels: map[Level]string{
		TRACE: "\033[7;95mTRACE\033[0m ",
		DEBUG: "\033[7;94mDEBUG\033[0m ",
		INFO:  "\033[7;92mINFO \033[0m ",
		WARN:  "\033[7;93mWARN \033[0m ",
		ERROR: "\033[7;91mERROR\033[0m ",
		FATAL: "\033[7;35mFATAL\033[0m ",
	},
	AltLevel: Pair{"\033[7", "\033[0m "},
	Category: Pair{"\033[2m[\033[22;1;92m", "\033[0;2m]:\033[0m "},
	Field:    Pair{"\n    \033[36m", "\033[0m="},
}
