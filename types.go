package log

import "io"

// Level описывает уровень записи лога.
type Level int8

// Предопределенные уровни сообщений для записи в лог.
const (
	TRACE Level = (iota - 2) * (1 << 5) // -64
	DEBUG                               // -32
	INFO                                // 0
	WARN                                // 32
	ERROR                               // 64
	FATAL                               // 96
)

// String возвращает название группы уровня события: "TRACE", "DEBUG", "INFO",
// "WARN", "ERROR" и "FATAL".
func (l Level) String() string {
	switch l & -32 {
	case -32:
		return "DEBUG"
	case 0:
		return "INFO"
	case 32:
		return "WARN"
	case 64:
		return "ERROR"
	case 96:
		return "FATAL"
	default:
		return "TRACE"
	}
}

// Handler описывает интерфейс для поддержки вывода записей в лог.
// При выводе каждой записи лога, вызывается метод Log, в который передаются
// все параметры. Форматирование записи лога и его сохранение или вывод целиком
// лежат на обработчике.
type Handler interface {
	Log(lvl Level, category, msg string, fields ...interface{}) error
}

// StreamFormatter описывает интерфейс для форматированного вывода лога.
type StreamFormatter interface {
	Log(w io.Writer, lvl Level, category, msg string, fields ...interface{}) error
}

// Fields задает синоним для списка именованных полей.
type Fields = map[string]interface{}

// Pair описывает префикс и суффикс для опции. Используется при задании
// префиксов и суффиксов настроек вывода консольного лога.
type Pair = [2]string
