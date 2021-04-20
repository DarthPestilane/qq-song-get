package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// colors
const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

var _ logrus.Formatter = &Formatter{}

// Formatter formats logs
type Formatter struct {
	// DisableColor if turns to true, formatter won't add color for the log levels
	DisableColor bool
}

// Format implements logurs.Formatter
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	ts := entry.Time.Format(time.RFC3339)
	msg := entry.Message

	levelTxt := strings.ToUpper(entry.Level.String()[:4])
	if !f.DisableColor {
		// print color
		var levelColor int
		switch entry.Level {
		case logrus.DebugLevel, logrus.TraceLevel:
			levelColor = gray
		case logrus.WarnLevel:
			levelColor = yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = red
		default:
			levelColor = blue
		}
		levelTxt = fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, levelTxt)
	}

	return []byte(fmt.Sprintf("[%s]%s: %s\n", ts, levelTxt, msg)), nil
}
