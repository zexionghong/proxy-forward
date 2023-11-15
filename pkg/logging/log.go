package logging

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type Level int

var (
	Log *logrus.Logger
)

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

// Setup initialize the log instance
func init() {
	Log = logrus.New()
	Log.Out = os.Stdout
	//Log.SetFormatter(&logrus.TextFormatter{
	//	TimestampFormat: "2006-01-02 15:04:05",
	//	FullTimestamp:   true,
	//})
	Log.SetFormatter(&MyFormatter{})
}
