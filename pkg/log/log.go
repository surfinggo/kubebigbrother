package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Logger *logrus.Logger
)

func init() {
	Logger = logrus.New()
	l, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		l = logrus.InfoLevel
	}
	Logger.Level = l
	Logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
}

func Trace(a ...interface{})                 { Logger.Trace(a...) }
func Debug(a ...interface{})                 { Logger.Debug(a...) }
func Info(a ...interface{})                  { Logger.Info(a...) }
func Warn(a ...interface{})                  { Logger.Warn(a...) }
func Error(a ...interface{})                 { Logger.Error(a...) }
func Fatal(a ...interface{})                 { Logger.Fatal(a...) }
func Tracef(format string, a ...interface{}) { Logger.Tracef(format, a...) }
func Debugf(format string, a ...interface{}) { Logger.Debugf(format, a...) }
func Infof(format string, a ...interface{})  { Logger.Infof(format, a...) }
func Warnf(format string, a ...interface{})  { Logger.Warnf(format, a...) }
func Errorf(format string, a ...interface{}) { Logger.Errorf(format, a...) }
func Fatalf(format string, a ...interface{}) { Logger.Fatalf(format, a...) }
