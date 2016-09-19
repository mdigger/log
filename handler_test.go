package log

import (
	"errors"
	"testing"
)

type BadHanlder struct{}

func (h *BadHanlder) Handle(e *Entry) error {
	return errors.New("bad handler")
}

func TestBadHandler(t *testing.T) {
	New(new(BadHanlder)).Info("message")
}
