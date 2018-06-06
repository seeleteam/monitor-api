/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package config

import (
	"github.com/sirupsen/logrus"
)

var (
	// APPName monitor app name
	APPName = "monitor-api"

	// VERSION represent seele monitor api version.
	VERSION = "0.1.0"

	// ShardMap shard:<websocket url>
	ShardMap map[string]string
)

const (

	// DEV is for develop
	DEV = "dev"

	// RELEASE is for production
	RELEASE = "release"

	// DefaultTimeUnitSecond default time unit for config
	DefaultTimeUnitSecond = "s"

	// DefaultTimeUnitMilliSecond default time unit for config
	DefaultTimeUnitMilliSecond = "ms"

	// DefaultTimeUnit default time unit for config
	DefaultTimeUnit = DefaultTimeUnitSecond
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel logrus.Level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

// LogLevelMap log level map for logrus
var LogLevelMap = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
}
