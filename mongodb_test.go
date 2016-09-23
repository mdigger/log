package log_test

import (
	"os"

	"github.com/mdigger/log"
	"gopkg.in/mgo.v2"
)

type MongoDBHandler struct {
	log.Level
	WithSource bool
	session    *mgo.Session
}

func NewMongoDBHandler(url string) (*MongoDBHandler, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	if err := session.Ping(); err != nil {
		return nil, err
	}
	return &MongoDBHandler{session: session}, nil
}

func (h *MongoDBHandler) Close() {
	if h.session != nil {
		h.session.Close()
	}
}

// Handle implements the log.Handler interface.
func (h *MongoDBHandler) Handle(e *log.Entry) error {
	// check that the handler supports writing to log messages of this level
	if e.Level < h.Level {
		return nil
	}
	// if you need information about the source file, get it
	if h.WithSource && e.Source == nil {
		e.Source = log.NewSource(5) // use the level of call stack
	}
	session := h.session.Copy()
	err := session.DB("").C("log").Insert(e)
	session.Close()
	return err
}

func ExampleHandler() {
	mdbHandler, err := NewMongoDBHandler("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	defer mdbHandler.Close()
	mdbHandler.WithSource = true

	log := log.New(
		log.NewConsole(os.Stderr, log.LstdFlags), // console log
		mdbHandler, // mongoDB log
	)
	log.Info("log started")
}
