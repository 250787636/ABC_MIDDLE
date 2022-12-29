package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// goland:noinspection ALL
var (
	Trace = logrus.Trace
	Debug = logrus.Debug
	Info  = logrus.Info
	Warn  = logrus.Warn
	Error = logrus.Error
	Fatal = logrus.Fatal
	Panic = logrus.Panic

	Tracef = logrus.Tracef
	Debugf = logrus.Debugf
	Infof  = logrus.Infof
	Warnf  = logrus.Warnf
	Errorf = logrus.Errorf
	Fatalf = logrus.Fatalf
	Panicf = logrus.Panicf

	printf = logrus.Printf

	DefalutLogger = logrus.StandardLogger()
	FileWriter    *lumberjack.Logger
)

// 配置logger
func init() {
	_ = os.MkdirAll("logs", 0777)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: false,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		QuoteEmptyFields:          true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			file = fmt.Sprintf("[ %s:%d ]", filepath.Base(frame.File), frame.Line)
			return
		},
	})

	logrus.SetReportCaller(true)
	FileWriter = &lumberjack.Logger{
		Filename:   "logs/run.log",
		MaxSize:    100,
		MaxAge:     1,
		MaxBackups: 31,
		Compress:   false,
	}

	logrus.SetOutput(io.MultiWriter(os.Stdout, FileWriter))
}
