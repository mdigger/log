package log

// Level задает уровень записи лога.
type Level int8

// Предопределенные уровни сообщений для записи в лог.
const (
	TRACE Level = (iota - 2) << 5 // -64
	DEBUG                         // -32
	INFO                          // 0
	WARN                          // 32
	ERROR                         // 64
	FATAL                         // 96
)

// String возвращает название группы уровней записи лога:
//  ""      [-128...-65]
//  "TRACE" [-64...-33]
//  "DEBUG" [-32...-1]
//  "INFO"  [0...31]
//  "WARN"  [32...63]
//  "ERROR" [64...95]
//  "FATAL" [96...127]
func (l Level) String() string {
	switch l & -32 {
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	default:
		return ""
	}
}
