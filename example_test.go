package log_test

import (
	"os"

	"github.com/mdigger/log"
)

func Example_console() {
	clog := log.NewConsoleHandler(os.Stdout, log.Lshortfile) // новый лог для вывода в консоль
	clog.Padding = 16                                        // длина сообщения для выравнивания
	logger := log.New(clog)

	logger.Info("info", "test")     // информационное сообщение
	logger.Infof("info %v", "test") // информационное сообщение

	logger.Error("error", "test")     // ошибка
	logger.Errorf("error %v", "test") // ошибка

	logger.Debug("debug", "test")     // отладочная информация (не выводится)
	logger.Debugf("debug %v", "test") // отладочная информация (не выводится)

	// информационное сообщение с дополнительными параметрами
	logger.WithField("key", "value").Info("test")
	logger.WithFields(log.Fields{
		"key":  "value",
		"key2": "value2",
	}).Info("test")

	// добавляем в параметры имя файла и номер строки с исходным кодом
	logger.WithSource(0, false).Info("info", "test")
	logger.WithField("key", "value").WithSource(0, false).Info("test")

	// Output:
	// example_test.go:14: info test
	// example_test.go:15: info test
	// example_test.go:17: error test
	// example_test.go:18: error test
	// example_test.go:24: test             key=value
	// example_test.go:28: test             key=value key2=value2
	// example_test.go:31: info test        source=example_test.go:31
	// example_test.go:32: test             key=value source=example_test.go:32
}

func Example_json() {
	logger := log.New(log.NewJSONHandler(os.Stdout, 0))

	logger.Info("info", "test")     // информационное сообщение
	logger.Infof("info %v", "test") // информационное сообщение

	logger.Error("error", "test")     // ошибка
	logger.Errorf("error %v", "test") // ошибка

	logger.Debug("debug", "test")     // отладочная информация (не выводится)
	logger.Debugf("debug %v", "test") // отладочная информация (не выводится)

	// информационное сообщение с дополнительными параметрами
	logger.WithField("key", "value").Info("test")
	logger.WithFields(log.Fields{
		"key":  "value",
		"key2": "value2",
	}).Info("test")

	// добавляем в параметры имя файла и номер строки с исходным кодом
	logger.WithSource(0, false).Info("info", "test")
	logger.WithField("key", "value").WithSource(0, false).Info("test")

	// Output:
	// {"level":"info","msg":"info test"}
	// {"level":"info","msg":"info test"}
	// {"level":"error","msg":"error test"}
	// {"level":"error","msg":"error test"}
	// {"level":"info","msg":"test","fields":{"key":"value"}}
	// {"level":"info","msg":"test","fields":{"key":"value","key2":"value2"}}
	// {"level":"info","msg":"info test","fields":{"source":"example_test.go:65"}}
	// {"level":"info","msg":"test","fields":{"key":"value","source":"example_test.go:66"}}
}

func Example_mixed() {
	// новый лог для вывода в консоль
	clog := log.NewConsoleHandler(os.Stdout, log.Lshortfile)
	clog.Padding = 16 // длина сообщения для выравнивания
	// новый лог для вывода в консоль
	json := log.NewJSONHandler(os.Stdout, 0)
	json.SetFlags(0) // сбрасываем флаги
	// вывод сразу в несколько логов в разных форматах
	logger := log.New(clog, json)
	logger.Info("info") // информационное сообщение

	// Output:
	// example_test.go:88: info
	// {"level":"info","msg":"info"}
}
