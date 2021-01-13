package log

import (
	"encoding/json"
	"errors"
)

var (
	// LogLevelMapping is a mapping for LogLevel enum
	LoglevelMapping = map[string]LogLevel{
		ERROR.String():   ERROR,
		WARNING.String(): WARNING,
		INFO.String():    INFO,
		DEBUG.String():   DEBUG,
		SILENT.String():  SILENT,
		FATAL.String():   FATAL,
	}
)

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	SILENT
	FATAL
)

type LogLevel int

// UnmarshalYAML unserialize LogLevel with yaml
func (l *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tp string
	unmarshal(&tp)
	level, exist := LoglevelMapping[tp]
	if !exist {
		return errors.New("invalid mode")
	}
	*l = level
	return nil
}

// UnmarshalJSON unserialize LogLevel with json
func (l *LogLevel) UnmarshalJSON(data []byte) error {
	var tp string
	json.Unmarshal(data, &tp)
	level, exist := LoglevelMapping[tp]
	if !exist {
		return errors.New("invalid mode")
	}
	*l = level
	return nil
}

// MarshalJSON serialize LogLevel with json
func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

// MarshalYAML serialize LogLevel with yaml
func (l LogLevel) MarshalYAML() (interface{}, error) {
	return l.String(), nil
}

func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "INFO"
	case WARNING:
		return "WARN"
	case ERROR:
		return "ERROR"
	case DEBUG:
		return "DEBUG"
	case SILENT:
		return "SILENT"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
