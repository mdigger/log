package log

// Fields represents a map of Context level data used for structured logging.
type Fields map[string]interface{}

// WithFields returns a new Fields with added fields.
func (f Fields) WithFields(fields Fields) Fields {
	if fields == nil {
		return f
	}
	for key, value := range f {
		if _, ok := fields[key]; !ok {
			fields[key] = value
		}
	}
	return fields
}

// WithField returns a new Fields with added the named field.
func (f Fields) WithField(name string, value interface{}) Fields {
	fields := make(Fields, len(f)+1)
	fields[name] = value
	return f.WithFields(fields)
}

// WithError returns a new Fields with added field "error" contains the error
// description.
func (f Fields) WithError(err error) Fields {
	if err == nil {
		return f
	}
	return f.WithField("error", err.Error())
}

// WithSource return new Fields with added information about the file name and
// line number of the source code. Calldepth is the count of the number of
// frames to skip when computing the file name and line number. A value of 0
// will print the details for the caller.
func (f Fields) WithSource(calldepth int) Fields {
	return f.WithField("source", NewSource(calldepth+1))
}
