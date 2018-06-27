/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package core

import (
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/seeleteam/monitor-api/config"
	"github.com/seeleteam/monitor-api/core/logs"
)

const (
	defaultHammerTime = 10 // if use graceful shutdown, this will effect
)

// MonitorServer monitor server config
type MonitorServer struct {
	Server   *http.Server
	CertFile string
	KeyFile  string
	G        *errgroup.Group
}

// GetServer config and return the MonitorServer
func GetServer(g *errgroup.Group, handlerConfig ...*EngineConfig) (slServer *MonitorServer) {
	var currentEngineConfig *EngineConfig

	currentServerConfig := config.SeeleConfig.ServerConfig
	if currentServerConfig == nil {
		logs.Fatal("error nil currentServerConfig")
	}

	defaultEngineConfig := currentServerConfig.EngineConfig

	if len(handlerConfig) == 0 {
		// use default config
		if defaultEngineConfig == nil {
			logs.Fatal("error nil defaultEngineConfig")
		}

		currentEngineConfig = &EngineConfig{
			DisableConsoleColor: defaultEngineConfig.DisableConsoleColor,
			WriteLog:            defaultEngineConfig.WriteLog,
			LogFile:             defaultEngineConfig.LogFile,
			TempFolder:          defaultEngineConfig.TempFolder,
			LimitConnections:    defaultEngineConfig.LimitConnection,
		}
	} else {
		// use yourself define config
		// if writeLog = true(write logs to file), set disableConsoleColor = true
		writeLog := handlerConfig[0].WriteLog
		logFile := handlerConfig[0].LogFile
		tempFloder := handlerConfig[0].TempFolder
		if tempFloder == "" {
			tempFloder = defaultEngineConfig.TempFolder
		}

		s, err := os.Stat(tempFloder)
		if err != nil {
			logs.Fatal("error tempFloder %v", err)
		}
		if !s.IsDir() {
			logs.Fatal("tempFloder %v is not dir", tempFloder)
		}

		disableConsoleColor := handlerConfig[0].DisableConsoleColor
		if writeLog {
			if len(logFile) == 0 {
				logFile = currentEngineConfig.LogFile
			}
			disableConsoleColor = true
		}

		limitConnections := handlerConfig[0].LimitConnections
		if limitConnections < 0 {
			limitConnections = 0
		}

		currentEngineConfig = &EngineConfig{
			DisableConsoleColor: disableConsoleColor,
			WriteLog:            writeLog,
			TempFolder:          tempFloder,
			LogFile:             logFile,
			LimitConnections:    limitConnections,
		}
	}

	return &MonitorServer{
		Server: &http.Server{
			Addr:           currentServerConfig.Addr,
			Handler:        currentEngineConfig.Init(),
			ReadTimeout:    currentServerConfig.ReadTimeout,
			WriteTimeout:   currentServerConfig.WriteTimeout,
			IdleTimeout:    currentServerConfig.IdleTimeout,
			MaxHeaderBytes: currentServerConfig.MaxHeaderBytes,
		},
		G: g,
	}
}

// NewServer create new MonitorServer with http.Server and errgroup.Group
func (sl *MonitorServer) NewServer(server *http.Server, g *errgroup.Group) {
	sl.Server = server
	sl.G = g
}

// NewServerTLS create new SSL MonitorServer with http.Server, errgroup.Group and TLS config file
func (sl *MonitorServer) NewServerTLS(server *http.Server, certFile string, keyFile string, g *errgroup.Group) {
	sl.Server = server
	sl.CertFile = certFile
	sl.KeyFile = keyFile
	sl.G = g
}

// RunServer run our server in a goroutine so that it doesn't block.
func (sl *MonitorServer) RunServer() {
	sl.G.Go(func() error {
		return http.ListenAndServe(sl.Server.Addr, sl.Server.Handler)
	})
}

// RunServerTLS run our server with tls in a goroutine so that it doesn't block.
func (sl *MonitorServer) RunServerTLS() {
	sl.G.Go(func() error {
		return http.ListenAndServeTLS(sl.Server.Addr, sl.CertFile, sl.KeyFile, sl.Server.Handler)
	})
}
