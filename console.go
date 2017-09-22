package log

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

// Console отвечает за форматирование лога в текстовом консольном виде.
// В нем задается формат представления временной метки и префиксы/суффиксы,
// используемые при выводе элементов лога.
type Console struct {
	TimeFormat string           // формат вывода даты и времени
	UTC        bool             // вывод даты и времени в UTC
	Time       Pair             // для обрамления временной метки
	Levels     map[Level]string // переопределение строк для вывода уровня
	AltLevel   Pair             // для альтернативных уровней
	Category   Pair             // для обрамления категории
	Field      Pair             // перед названием и между названием и значением
}

var _ StreamFormatter = new(Console) // проверяем поддержу интерфейса

// Log отвечает за форматирование и вывод записи лога.
func (f Console) Log(w io.Writer, lvl Level, category, msg string,
	fields ...interface{}) error {
	var buf = buffers.Get().([]byte) // получаем новый буфер
	buf = buf[:0]                    // сбрасываем возможные предыдущие значения
	// выводим дату и время, если задан формат
	if f.TimeFormat != "" {
		var ts = time.Now()
		if f.UTC {
			ts = ts.UTC()
		}
		if f.Time[0] != "" {
			buf = append(buf, f.Time[0]...)
		}
		buf = ts.AppendFormat(buf, f.TimeFormat)
		if f.Time[1] != "" {
			buf = append(buf, f.Time[1]...)
		} else {
			buf = append(buf, ' ')
		}
	}
	// уровень записи
repeatLevel:
	switch level, ok := f.Levels[lvl]; {
	case ok:
		if level != "" {
			buf = append(buf, level...)
		}
	case lvl != -64 && lvl != -32 && lvl != 0 && lvl != 32 && lvl != 64 && lvl != 96:
		// нормализуем уровень до идентификаторов интервалов
		if lvl > -32 {
			lvl &= -32
		} else {
			lvl = -64 // TRACE
		}
		goto repeatLevel
	default:
		if f.AltLevel[0] != "" {
			buf = append(buf, f.AltLevel[0]...)
		}
		buf = append(buf, lvl.String()...)
		if f.AltLevel[1] != "" {
			buf = append(buf, f.AltLevel[1]...)
		} else {
			buf = append(buf, ' ')
		}
	}
	// категория
	if category != "" {
		if f.Category[0] != "" {
			buf = append(buf, f.Category[0]...)
		} else {
			buf = append(buf, '[')
		}
		buf = append(buf, category...)
		if f.Category[1] != "" {
			buf = append(buf, f.Category[1]...)
		} else {
			buf = append(buf, "]: "...)
		}
	}
	// основной текст
	if msg != "" {
		buf = append(buf, msg...)
	}
	// дополнительные поля
	if f.Field[0] == "" {
		f.Field[0] = " "
	}
	if f.Field[1] == "" {
		f.Field[1] = "="
	}
	switch len(fields) {
	case 0: // нет дополнительных полей
		break
	case 1: // дополнительные поля представлены одним элементом
		if list, ok := fields[0].(map[string]interface{}); ok {
			for name, value := range list {
				buf = append(buf, f.Field[0]...)
				buf = append(buf, name...)
				buf = append(buf, f.Field[1]...)
				strValue(&buf, value)
			}
		}
	default:
		for i, field := range fields {
			if i%2 == 0 {
				buf = append(buf, f.Field[0]...)
			} else {
				buf = append(buf, f.Field[1]...)
			}
			strValue(&buf, field)
		}
	}
	buf = append(buf, '\n')
	_, err := w.Write(buf)
	buffers.Put(buf)
	return err
}

var buffers = sync.Pool{New: func() interface{} { return []byte{} }}

func strValue(buf *[]byte, value interface{}) {
	switch value := value.(type) {
	case string:
		*buf = append(*buf, value...)
	case []byte:
		*buf = append(*buf, value...)
	case error:
		*buf = append(*buf, value.Error()...)
	case fmt.Stringer:
		*buf = append(*buf, value.String()...)
	case bool:
		*buf = strconv.AppendBool(*buf, value)
	case int:
		*buf = strconv.AppendInt(*buf, int64(value), 10)
	case int8:
		*buf = strconv.AppendInt(*buf, int64(value), 10)
	case int16:
		*buf = strconv.AppendInt(*buf, int64(value), 10)
	case int32:
		*buf = strconv.AppendInt(*buf, int64(value), 10)
	case int64:
		*buf = strconv.AppendInt(*buf, value, 10)
	case uint:
		*buf = strconv.AppendUint(*buf, uint64(value), 10)
	case uint8:
		*buf = strconv.AppendUint(*buf, uint64(value), 10)
	case uint16:
		*buf = strconv.AppendUint(*buf, uint64(value), 10)
	case uint32:
		*buf = strconv.AppendUint(*buf, uint64(value), 10)
	case uint64:
		*buf = strconv.AppendUint(*buf, value, 10)
	case float32:
		*buf = strconv.AppendFloat(*buf, float64(value), 'g', -1, 32)
	case float64:
		*buf = strconv.AppendFloat(*buf, value, 'g', -1, 64)
	default:
		*buf = append(*buf, fmt.Sprint(value)...)
	}
}
