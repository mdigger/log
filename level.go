package log

// LogLevel is the log severity level.
type LogLevel int8

// String converts a severity level to a string.
func (l LogLevel) String() string {
	switch {
	case l < 0:
		return "debug"
	case l > 0:
		return "error"
	default:
		return "info"
	}
}

// MarshalText return the level string for text based encoding.
func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// Predefined log severity levels.
const (
	LevelDebug LogLevel = -1 // debug level
	LevelInfo  LogLevel = 0  // default severity level
	LevelError LogLevel = 1  // error level
)
