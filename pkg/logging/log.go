package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Level int

var (
	Log *logrus.Logger
)

// Setup initialize the log instance
func init() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}
