/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package config

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/seeleteam/monitor-api/core/config"
	"github.com/seeleteam/monitor-api/core/utils"
)

// Config is the main struct for SeeleConfig
type Config struct {
	AppName      string //Application name
	RunMode      string //Running Mode: dev | release
	RecoverFunc  func(*gin.Context)
	RecoverPanic bool
	ServerConfig *ServerConfig // server config
	ServerName   string
}

// ServerConfig define the parameters for running an HTTP server.
type ServerConfig struct {
	Addr     string // format must be like ip:port, :port, domain or domain:port
	ErrorLog *log.Logger
	// Server config

	Handler           http.Handler // need config the routers and filter
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	LogLevel          logrus.Level // log level
	TLSConfig         *tls.Config
	onShutdown        []func()
	WriteTimeout      time.Duration

	// Engine config
	EngineConfig *EngineConfig // config for gin engine

	// WebSocket config
	EnableWebSocket bool
	WebSocketConfig *WebSocketConfig

	// RPC config
	EnableRPC bool
	RPCConfig *RPCConfig
}

// EngineConfig for the router config
type EngineConfig struct {
	DisableConsoleColor bool
	LimitConnection     int
	LogFile             string
	WriteLog            bool
}

// RPCConfig for the rpc config
type RPCConfig struct {
	Debug  bool
	Scheme string
	URL    string
}

// WebSocketConfig is the base webSocket config
type WebSocketConfig struct {
	DelayReConnTime       time.Duration // delay recon time when web socket error occur
	DelaySendTime         time.Duration // delay recon and resend msg to monitor when rpc server occur error
	ReportErrorAfterTimes int           // report the error when error occur over the special times

	WsFullEventTickerTime        time.Duration // send msg with ticker
	WsLatestBlockEventTickerTime time.Duration // send msg with ticker
	WsPass                       string        // Password to authorize access to the monitoring page
	WsRouter                     string        // full path host:port/ws and the WsRouter is /ws
	WsURL                        string        // host:port
}

var (
	// AppConfig is the instance of Configure, store the config information from file
	AppConfig *MonitorAppConfig
	// AppPath is the absolute path to the app
	AppPath string

	// appConfigPath is the path to the config files
	appConfigPath string
	// appConfigProvider is the provider for the config, default is ini
	appConfigProvider = "ini"

	// SeeleConfig is the default config for Application
	SeeleConfig *Config
)

// Init init the config
func Init(configFile string) {
	SeeleConfig = newSeeleConfig()
	var err error
	if err = parseConfig(configFile); err != nil {
		panic(err)
	}
}

func newSeeleConfig() *Config {
	defaultAddr := ":9999"

	defaultDelayReConnTime := 5 * time.Second
	defaultDelaySendTime := 5 * time.Second

	defaultLogLevel := InfoLevel

	defaultWsFullEventTickerTime := 10 * time.Second
	defaultWsLatestEventTickerTime := 5 * time.Second
	defaultWsRouter := "/api"

	return &Config{
		AppName:      APPName,
		RecoverPanic: true,
		RunMode:      DEV,
		ServerName:   "monitor-api:" + VERSION,
		ServerConfig: &ServerConfig{
			Addr:              defaultAddr,
			IdleTimeout:       0,
			LogLevel:          defaultLogLevel,
			MaxHeaderBytes:    1 << 20, //1MB
			ReadTimeout:       300 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
			WriteTimeout:      120 * time.Second,
			EngineConfig: &EngineConfig{
				DisableConsoleColor: false,
				WriteLog:            false,
				LogFile:             APPName + ".log",
				LimitConnection:     0, // no limits the conn per timeUnit
			},
			EnableWebSocket: false,
			WebSocketConfig: &WebSocketConfig{
				WsURL:                        defaultAddr,
				WsRouter:                     defaultWsRouter,
				WsFullEventTickerTime:        defaultWsFullEventTickerTime,
				WsLatestBlockEventTickerTime: defaultWsLatestEventTickerTime,
				DelayReConnTime:              defaultDelayReConnTime,
				DelaySendTime:                defaultDelaySendTime,
				ReportErrorAfterTimes:        10,
				WsPass:                       "",
			},
			EnableRPC: false,
			RPCConfig: &RPCConfig{
				URL:    "", // rpc url
				Scheme: "tcp",
				Debug:  false,
			},
		},
	}
}

// parseConfig only support ini
func parseConfig(appConfigPath string) (err error) {
	AppConfig, err = newAppConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return err
	}
	return assignConfig(AppConfig)
}

// assignConfig assign the config
func assignConfig(ac config.Configure) error {
	for _, i := range []interface{}{SeeleConfig, &SeeleConfig.ServerConfig} {
		assignSingleConfig(i, ac)
	}
	// set the run mode first, env set is the highest priority
	if envRunMode := os.Getenv("MONITOR_API_RUNMODE"); envRunMode != "" {
		SeeleConfig.RunMode = envRunMode
	} else if runMode := ac.String("RunMode"); runMode != "" {
		SeeleConfig.RunMode = runMode
	} else {
		defaultEnv, err := ac.GetSection("default")
		if err != nil {
			SeeleConfig.RunMode = DEV
		} else {
			SeeleConfig.RunMode = defaultEnv["run_mode"]
		}
	}

	// APPName load from default, default section, real section
	currentAppName := APPName
	if defaultSection, err := ac.GetSection("default"); err == nil {
		if len(defaultSection["app_name"]) != 0 {
			currentAppName = defaultSection["app_name"]
		}
	}

	// first use default section, and use real mode to override
	currentSection, err := ac.GetSection(SeeleConfig.RunMode)
	if err != nil {
		SeeleConfig.RunMode = DEV
		currentSection, _ = ac.GetSection(SeeleConfig.RunMode)
	}
	// real APPName
	if currentSection["app_name"] != "" {
		currentAppName = currentSection["app_name"]
	}
	SeeleConfig.AppName = currentAppName
	APPName = currentAppName

	currentServerConfig := SeeleConfig.ServerConfig
	if currentSection["addr"] != "" {
		currentServerConfig.Addr = currentSection["addr"]
	}
	if currentSection["readtimeout"] != "" {
		currentReadTimeout, err := time.ParseDuration(currentSection["readtimeout"] + DefaultTimeUnit)
		if err == nil {
			currentServerConfig.ReadTimeout = currentReadTimeout
		}
	}
	if currentSection["readheadertimeout "] != "" {
		currentReadHeaderTimeout, err := time.ParseDuration(currentSection["readheadertimeout"] + DefaultTimeUnit)
		if err == nil {
			currentServerConfig.ReadHeaderTimeout = currentReadHeaderTimeout
		}
	}
	if currentSection["writetimeout "] != "" {
		currentWriteTimeout, err := time.ParseDuration(currentSection["writetimeout"] + DefaultTimeUnit)
		if err == nil {
			currentServerConfig.WriteTimeout = currentWriteTimeout
		}
	}
	if currentSection["idletimeout "] != "" {
		currentIdleTimeout, err := time.ParseDuration(currentSection["idletimeout"] + DefaultTimeUnit)
		if err == nil {
			currentServerConfig.IdleTimeout = currentIdleTimeout
		}
	}
	if currentSection["maxheaderbytes"] != "" {
		currentMaxHeaderBytes, err := strconv.Atoi(currentSection["maxheaderbytes"])
		if err == nil {
			currentServerConfig.MaxHeaderBytes = currentMaxHeaderBytes
		}
	}

	if currentSection["loglevel"] != "" {
		currentLogLevel, ok := LogLevelMap[currentSection["loglevel"]]
		if !ok {
			currentLogLevel = InfoLevel
		}
		currentServerConfig.LogLevel = currentLogLevel
	}

	if currentSection["enablewebsocket"] != "" {
		currentEnableWebSocket, err := strconv.ParseBool(currentSection["enablewebsocket"])
		if err == nil {
			currentServerConfig.EnableWebSocket = currentEnableWebSocket

			// web socket config deal
			currentWebSocketConfig := currentServerConfig.WebSocketConfig
			if currentWebSocketConfig != nil {
				if currentSection["wsurl"] != "" {
					currentWebSocketURL := currentSection["wsurl"]
					if err == nil {
						currentWebSocketConfig.WsURL = currentWebSocketURL
					}
				}
				if currentSection["wsrouter"] != "" {
					currentWsRouter := currentSection["wsrouter"]
					if err == nil {
						currentWebSocketConfig.WsRouter = currentWsRouter
					}
				}
				if currentSection["wsfulleventtickertime"] != "" {
					currentWsFullEventTickerTime, err := time.ParseDuration(currentSection["wsfulleventtickertime"] + DefaultTimeUnit)
					if err == nil {
						currentWebSocketConfig.WsFullEventTickerTime = currentWsFullEventTickerTime
					}
				}
				if currentSection["wslatestblockeventtickertime"] != "" {
					currentWsLatestBlockEventTickerTime, err := time.ParseDuration(currentSection["wslatestblockeventtickertime"] + DefaultTimeUnit)
					if err == nil {
						currentWebSocketConfig.WsLatestBlockEventTickerTime = currentWsLatestBlockEventTickerTime
					}
				}
				if currentSection["delayreconntime"] != "" {
					currentDelayReConnTime, err := time.ParseDuration(currentSection["delayreconntime"] + DefaultTimeUnit)
					if err == nil {
						currentWebSocketConfig.DelayReConnTime = currentDelayReConnTime
					}
				}
				if currentSection["delaysendtime"] != "" {
					currentDelaySendTime, err := time.ParseDuration(currentSection["delaysendtime"] + DefaultTimeUnit)
					if err == nil {
						currentWebSocketConfig.DelaySendTime = currentDelaySendTime
					}
				}
				if currentSection["reporterroraftertimes"] != "" {
					currentReportErrorAfterTimes, err := strconv.Atoi(currentSection["reporterroraftertimes"])
					if err == nil {
						currentWebSocketConfig.ReportErrorAfterTimes = currentReportErrorAfterTimes
					}
				}

				currentServerConfig.WebSocketConfig = currentWebSocketConfig
			}
		}
	}
	if currentSection["enablerpc"] != "" {
		currentEnableRPC, err := strconv.ParseBool(currentSection["enablerpc"])
		if err == nil {
			currentServerConfig.EnableRPC = currentEnableRPC

			// rpc config deal
			currentRPCConfig := currentServerConfig.RPCConfig
			if currentRPCConfig != nil {
				if currentSection["rpcurl"] != "" {
					currentRPCURL := currentSection["rpcurl"]
					if err == nil {
						currentRPCConfig.URL = currentRPCURL
					}
				}
				currentServerConfig.RPCConfig = currentRPCConfig
			}
		}
	}

	// engine config deal
	currentEngineConfig := currentServerConfig.EngineConfig
	if currentEngineConfig != nil {
		if currentSection["limitconnection"] != "" {
			currentLimitConnection, err := strconv.Atoi(currentSection["limitconnection"])
			if err == nil {
				currentEngineConfig.LimitConnection = currentLimitConnection
			}
		}
		if currentSection["disableconsolecolor"] != "" {
			currentDisableConsoleColor, err := strconv.ParseBool(currentSection["disableconsolecolor"])
			if err == nil {
				currentEngineConfig.DisableConsoleColor = currentDisableConsoleColor
			}
		}
		if currentSection["writelog"] != "" {
			currentWriteLog, err := strconv.ParseBool(currentSection["writelog"])
			if err == nil {
				currentEngineConfig.WriteLog = currentWriteLog
				if currentSection["logfile"] != "" {
					currentEngineConfig.LogFile = currentSection["logfile"]
				}
			}
		}
		currentServerConfig.EngineConfig = currentEngineConfig
	}

	return nil
}

func assignSingleConfig(p interface{}, ac config.Configure) {
	pt := reflect.TypeOf(p)
	if pt.Kind() != reflect.Ptr {
		return
	}
	pt = pt.Elem()
	if pt.Kind() != reflect.Struct {
		return
	}
	pv := reflect.ValueOf(p).Elem()

	for i := 0; i < pt.NumField(); i++ {
		pf := pv.Field(i)
		if !pf.CanSet() {
			continue
		}
		name := pt.Field(i).Name
		switch pf.Kind() {
		case reflect.String:
			pf.SetString(ac.DefaultString(name, pf.String()))
		case reflect.Int, reflect.Int64:
			pf.SetInt(ac.DefaultInt64(name, pf.Int()))
		case reflect.Bool:
			pf.SetBool(ac.DefaultBool(name, pf.Bool()))
		case reflect.Struct:
		default:
			//do nothing here
		}
	}

}

// LoadAppConfig load a config file
func LoadAppConfig(adapterName, configPath string) error {
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	if !utils.FileExists(absConfigPath) {
		return fmt.Errorf("the target config file: %s don't exist", configPath)
	}

	appConfigPath = absConfigPath
	appConfigProvider = adapterName

	return parseConfig(appConfigPath)
}

// MonitorAppConfig monitor app config
type MonitorAppConfig struct {
	innerConfig config.Configure
}

func newAppConfig(appConfigProvider, appConfigPath string) (*MonitorAppConfig, error) {
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &MonitorAppConfig{ac}, nil
}

// Set set value for key
func (b *MonitorAppConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(SeeleConfig.RunMode+"::"+key, val); err != nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

// String get the string value with special key
func (b *MonitorAppConfig) String(key string) string {
	if v := b.innerConfig.String(SeeleConfig.RunMode + "::" + key); v != "" {
		return v
	}
	return b.innerConfig.String(key)
}

// Strings get the string array with special key
func (b *MonitorAppConfig) Strings(key string) []string {
	if v := b.innerConfig.Strings(SeeleConfig.RunMode + "::" + key); len(v) > 0 {
		return v
	}
	return b.innerConfig.Strings(key)
}

// Int get the int value with special key
func (b *MonitorAppConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(SeeleConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(key)
}

// Int64 get the int64 value with special key
func (b *MonitorAppConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(SeeleConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(key)
}

// Bool get the bool value with special key
func (b *MonitorAppConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(SeeleConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(key)
}

// Float get the float64 value with special key
func (b *MonitorAppConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(SeeleConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(key)
}

// DefaultString get the string value with special key, if value == "" return defaultVal
func (b *MonitorAppConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

// DefaultStrings get the string array with special key, if len(value) == 0 return defaultVal array
func (b *MonitorAppConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

// DefaultInt get the int value with special key, if value is not int return defaultVal
func (b *MonitorAppConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

// DefaultInt64 get the int64 value with special key, if value is not int64 return defaultVal
func (b *MonitorAppConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

// DefaultBool get the bool value with special key, if value is not bool return defaultVal
func (b *MonitorAppConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

// DefaultFloat get the float64 value with special key, if value is not float return defaultVal
func (b *MonitorAppConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

// DIY get the interface{} value with special key
func (b *MonitorAppConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

// GetSection get the section values with special section name
func (b *MonitorAppConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

// SaveConfigFile save the config into file
func (b *MonitorAppConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}
