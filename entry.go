package log

import "errors"

// Entry описывает частично заполненную дополнительными полями запись лога.
type Entry struct {
	logger *Logger
	fields Fields
}

// WithFields добавляет список дополнительных полей к записи.
func (c *Entry) WithFields(fields Fields) *Entry {
	if c.fields == nil {
		c.fields = fields
	} else {
		for name, value := range fields {
			c.fields[name] = value
		}
	}
	return c
}

// WithField добавляет именованный параметр к записи лога.
func (c *Entry) WithField(name string, value interface{}) *Entry {
	if c.fields == nil {
		c.fields = Fields{name: value}
	} else {
		c.fields[name] = value
	}
	return c
}

// WithError добавляет к дополнительным атрибутам не пустое значение ошибки.
func (c *Entry) WithError(err error) *Entry {
	if err != nil {
		return c.WithField("err", err)
	}
	return c
}

// WithSource добавляет в дополнительные атрибуты информацию об исходном файле.
func (c *Entry) WithSource() *Entry {
	return c.WithField("src", Source(0))
}

// Trace выводит необязательное отладочное сообщение в лог.
func (c *Entry) Trace(msg string) {
	c.logger.h.Log(TRACE, c.logger.name, msg, c.fields)
}

// Debug выводит отладочное сообщение в лог.
func (c *Entry) Debug(msg string) {
	c.logger.h.Log(DEBUG, c.logger.name, msg, c.fields)
}

// Info выводит информационное сообщение в лог.
func (c *Entry) Info(msg string) {
	c.logger.h.Log(INFO, c.logger.name, msg, c.fields)
}

// Warn выводит в лог предупреждение.
func (c *Entry) Warn(msg string) error {
	c.logger.h.Log(WARN, c.logger.name, msg, c.fields)
	return errors.New(msg)
}

// Error выводит в лок сообщение об ошибке.
func (c *Entry) Error(msg string) error {
	c.logger.h.Log(ERROR, c.logger.name, msg, c.fields)
	return errors.New(msg)
}

// Fatal выводит в лок высокоприоритетное сообщение об ошибке.
func (c *Entry) Fatal(msg string) error {
	c.logger.h.Log(FATAL, c.logger.name, msg, c.fields)
	return errors.New(msg)
}
