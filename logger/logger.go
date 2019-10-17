package logger

import (
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(defaultFormatter)
}

// SetFormatter override the default formatter set before
func SetFormatter(formatter logrus.Formatter) {
	logrus.SetFormatter(formatter)
}
