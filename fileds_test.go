package log

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestFields(t *testing.T) {
	var fields Fields
	if err := testFieldsLength(fields.WithField("name", "value"), 1); err != nil {
		t.Error(err)
	}
	if err := testFieldsLength(fields.WithFields(Fields{
		"name2": "value2",
		"time":  time.Now(),
	}), 2); err != nil {
		t.Error(err)
	}
	if err := testFieldsLength(fields.WithFields(Fields{
		"name2": "value2",
		"time":  time.Now(),
	}).WithField("name2", true), 2); err != nil {
		t.Error(err)
	}
	if err := testFieldsLength(fields.WithFields(Fields{
		"name2": "value2",
		"time":  time.Now(),
	}).WithField("name", "value"), 3); err != nil {
		t.Error(err)
	}
	if err := testFieldsLength(fields.WithSource(0), 1); err != nil {
		t.Error(err)
	}
	if err := testFieldsLength(fields.WithError(errors.New("error")), 1); err != nil {
		t.Error(err)
	}
	if err := testFieldsLength(fields.WithError(nil), 0); err != nil {
		t.Error(err)
	}
}

func testFieldsLength(fields Fields, length int) error {
	if l := len(fields); l != length {
		return fmt.Errorf("bad fields length: %v", l)
	}
	return nil
}
