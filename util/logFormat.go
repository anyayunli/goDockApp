package util

import (
	"bytes"
	"curbside-eta/pkg/color"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

// LogFormatter ...
type LogFormatter struct{}

// Format ...
func (f LogFormatter) Format(e *logrus.Entry) ([]byte, error) {
	_, fn, line, _ := runtime.Caller(7)
	if filepath.Base(fn) == "exported.go" {
		_, fn, line, _ = runtime.Caller(8)
	}
	var buffer bytes.Buffer
	var level string
	switch e.Level {
	case logrus.DebugLevel:
		level = color.Magenta("DEBU")
	case logrus.InfoLevel:
		level = color.Cyan("INFO")
	case logrus.WarnLevel:
		level = color.Yellow("WARN")
	case logrus.ErrorLevel:
		level = color.Red("ERRO")
	case logrus.FatalLevel:
		level = color.Red("FATA")
	case logrus.PanicLevel:
		level = color.Red("PANI")
	}
	buffer.WriteString(e.Time.Format("15:04:05"))
	buffer.WriteString(" ")
	buffer.WriteString(level)
	buffer.WriteString(" ")
	buffer.WriteString(color.Magenta("[" + filepath.Base(fn) + ":" + strconv.Itoa(line) + "]"))
	buffer.WriteString(" ")
	buffer.WriteString(e.Message)
	buffer.WriteString("\n")
	return buffer.Bytes(), nil
}
