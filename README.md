# A simple structured logging

[![GoDoc](https://godoc.org/github.com/mdigger/log?status.svg)](https://godoc.org/github.com/mdigger/log)
[![Build Status](https://travis-ci.org/mdigger/log.svg)](https://travis-ci.org/mdigger/log)
[![Coverage Status](https://coveralls.io/repos/github/mdigger/log/badge.svg?branch=master)](https://coveralls.io/github/mdigger/log?branch=master)

In general, this is another "bike" for logging with blackjack.

	package log

	import (
		"os"
		"time"

		"github.com/apex/log"
	)

	func open(filename string) (err error) {
		defer log.WithField("file", "~README.md").Trace("open").Stop(&err)
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 10)
		file.Close()
	}

	func main() {
		log.Info("info message")
		log.WithField("time", time.Now()).Debug("debug")

		err := open("README.md")
		if err != nil {
			// ...
		}
		// ,,,
	}
