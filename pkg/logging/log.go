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
func Setup() {
	Log = logrus.New()
	Log.Out = os.Stdout
	Log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
}
