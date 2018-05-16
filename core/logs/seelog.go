/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package logs

import (
	"fmt"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/seeleteam/monitor-api/config"
)

var logs *logrus.Logger

const (
	defaultLogPath = "monitor-api-logs"
	defaultLogFile = "monitor-api-all.logs"
)

// NewLogger create the logrus.Logger with special config
func NewLogger() *logrus.Logger {
	if logs != nil {
		return logs
	}

	logs = logrus.New()

	// get logLevel
	logLevel := config.SeeleConfig.ServerConfig.LogLevel
	logs.SetLevel(logLevel)

	writeLog := config.SeeleConfig.ServerConfig.EngineConfig.WriteLog
	if writeLog {
		_ = os.Mkdir(defaultLogPath, 0777)

		path := defaultLogPath + string(os.PathSeparator) + defaultLogFile
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),       // 24 hours
			rotatelogs.WithRotationTime(time.Duration(86400)*time.Second), // 1 days
		)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}

		logs.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.DebugLevel: writer,
				logrus.InfoLevel:  writer,
				logrus.WarnLevel:  writer,
				logrus.ErrorLevel: writer,
				logrus.FatalLevel: writer,
			},
			&logrus.JSONFormatter{},
		))

		pathMap := lfshook.PathMap{
			logrus.DebugLevel: fmt.Sprintf("%s/debug.log", defaultLogPath),
			logrus.InfoLevel:  fmt.Sprintf("%s/info.log", defaultLogPath),
			logrus.WarnLevel:  fmt.Sprintf("%s/warn.log", defaultLogPath),
			logrus.ErrorLevel: fmt.Sprintf("%s/error.log", defaultLogPath),
			logrus.FatalLevel: fmt.Sprintf("%s/fatal.log", defaultLogPath),
		}
		logs.AddHook(lfshook.NewHook(
			pathMap,
			&logrus.TextFormatter{},
		))
	}

	return logs
}

// GetLogger get the default logger
func GetLogger() *logrus.Logger {
	return logs
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

// Debug wrapper Debug logger
func Debug(f interface{}, args ...interface{}) {
	logs.Debug(formatLog(f, args...))
}

// Info wrapper Info logger
func Info(f interface{}, args ...interface{}) {
	logs.Info(formatLog(f, args...))
}

// Warn wrapper Warn logger
func Warn(f interface{}, args ...interface{}) {
	logs.Warn(formatLog(f, args...))
}

// Printf wrapper Printf logger
func Printf(f interface{}, args ...interface{}) {
	logs.Print(formatLog(f, args...))
}

// Panic wrapper Panic logger
func Panic(f interface{}, args ...interface{}) {
	logs.Panic(formatLog(f, args...))
}

// Fatal wrapper Fatal logger
func Fatal(f interface{}, args ...interface{}) {
	logs.Fatal(formatLog(f, args...))
}

// Error wrapper Error logger
func Error(f interface{}, args ...interface{}) {
	logs.Error(formatLog(f, args...))
}

// Debugln wrapper Debugln logger
func Debugln(v ...interface{}) {
	logs.Debugln(v...)
}

// Infoln wrapper Infoln logger
func Infoln(args ...interface{}) {
	logs.Infoln(args...)
}

// Warnln wrapper Warnln logger
func Warnln(args ...interface{}) {
	logs.Warnln(args...)
}

// Printfln wrapper Printfln logger
func Printfln(args ...interface{}) {
	logs.Println(args...)
}

// Panicln wrapper Panicln logger
func Panicln(args ...interface{}) {
	logs.Panicln(args...)
}

// Fatalln wrapper Fatalln logger
func Fatalln(args ...interface{}) {
	logs.Fatalln(args...)
}

// Errorln wrapper Errorln logger
func Errorln(args ...interface{}) {
	logs.Errorln(args...)
}
