package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type Formatter struct{}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("%s %s\n", timestamp, entry.Message)

	return []byte(msg), nil
}
