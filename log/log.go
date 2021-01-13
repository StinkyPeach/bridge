package log

import (
	"fmt"
	"os"

	"github.com/StinkyPeach/bridge/common/observable"

	log "github.com/sirupsen/logrus"
)

var (
	logCh  = make(chan interface{})
	source = observable.NewObservable(logCh)
	level  = INFO
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(new(Formatter))
}

type Event struct {
	LogLevel LogLevel
	Payload  string
}

type Event2 struct {
	LogLevel LogLevel
	Payload  []interface{}
}

func (e *Event) Type() string {
	return e.LogLevel.String()
}

func Info(format string, v ...interface{}) {
	event := newLog(INFO, format, v...)
	logCh <- event
	print(event)
}

func Warn(format string, v ...interface{}) {
	event := newLog(WARNING, format, v...)
	logCh <- event
	print(event)
}

func Error(format string, v ...interface{}) {
	event := newLog(ERROR, format, v...)
	logCh <- event
	print(event)
}

func Debug(format string, v ...interface{}) {
	event := newLog(DEBUG, format, v...)
	logCh <- event
	print(event)
}

func Fatal(format string, v ...interface{}) {
	event := newLog(FATAL, format, v...)
	logCh <- event
	print(event)
}

func Subscribe() observable.Subscription {
	sub, _ := source.Subscribe()
	return sub
}

func UnSubscribe(sub observable.Subscription) {
	source.UnSubscribe(sub)
}

func Level() LogLevel {
	return level
}

func SetLevel(newLevel LogLevel) {
	level = newLevel
}

func print(data *Event) {
	if data.LogLevel < level {
		return
	}

	switch data.LogLevel {
	case INFO:
		log.Infoln(fmt.Sprintf("[%s] ", data.LogLevel.String()) + data.Payload)
	case WARNING:
		log.Warnln(fmt.Sprintf("[%s] ", data.LogLevel.String()) + data.Payload)
	case ERROR:
		log.Errorln(fmt.Sprintf("[%s] ", data.LogLevel.String()) + data.Payload)
	case DEBUG:
		log.Debugln(fmt.Sprintf("[%s] ", data.LogLevel.String()) + data.Payload)
	case FATAL:
		log.Fatalln(fmt.Sprintf("[%s] ", data.LogLevel.String()) + data.Payload)
	}
}

func newLog(logLevel LogLevel, format string, v ...interface{}) *Event {
	return &Event{
		LogLevel: logLevel,
		Payload:  fmt.Sprintf(format, v...),
	}
}
