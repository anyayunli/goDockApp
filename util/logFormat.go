package util

import (
	"bytes"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

// NoColor ...
var NoColor = false

// Terminal styling constants
const (
	knrm = "\x1B[0m"
	kred = "\x1B[31m"
	kgrn = "\x1B[32m"
	kyel = "\x1B[33m"
	kblu = "\x1B[34m"
	kmag = "\x1B[35m"
	kcyn = "\x1B[36m"
	kwht = "\x1B[37m"
)

func colorStr(color string, val string) string {
	if NoColor {
		return val
	}
	return color + val + knrm
}

// White ...
func White(val string) string {
	return colorStr(kwht, val)
}

// Cyan ...
func Cyan(val string) string {
	return colorStr(kcyn, val)
}

// Red ...
func Red(val string) string {
	return colorStr(kred, val)
}

// Blue ...
func Blue(val string) string {
	return colorStr(kblu, val)
}

// Yellow ...
func Yellow(val string) string {
	return colorStr(kyel, val)
}

// Green ...
func Green(val string) string {
	return colorStr(kgrn, val)
}

// Magenta ...
func Magenta(val string) string {
	return colorStr(kmag, val)
}

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
		level = Magenta("DEBU")
	case logrus.InfoLevel:
		level = Cyan("INFO")
	case logrus.WarnLevel:
		level = Yellow("WARN")
	case logrus.ErrorLevel:
		level = Red("ERRO")
	case logrus.FatalLevel:
		level = Red("FATA")
	case logrus.PanicLevel:
		level = Red("PANI")
	}
	buffer.WriteString(e.Time.Format("15:04:05"))
	buffer.WriteString(" ")
	buffer.WriteString(level)
	buffer.WriteString(" ")
	buffer.WriteString(Magenta("[" + filepath.Base(fn) + ":" + strconv.Itoa(line) + "]"))
	buffer.WriteString(" ")
	buffer.WriteString(e.Message)
	buffer.WriteString("\n")
	return buffer.Bytes(), nil
}
