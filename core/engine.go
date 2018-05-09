/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package core

import (
	"net/http"

	"github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"

	"github.com/seeleteam/monitor-api/api/routers"
	"github.com/seeleteam/monitor-api/core/logs"
)

// EngineConfig gin engine config
type EngineConfig struct {
	DisableConsoleColor bool // disable the console color
	WriteLog            bool
	LogFile             string
	LimitConnections    int
	Routers             []gin.IRoutes
}

// initEngineConfig init engine config
func (config *EngineConfig) initEngineConfig() *gin.Engine {
	if config == nil {
		panic("engine config-watcher should not be nil")
	}

	if config.DisableConsoleColor {
		gin.DisableConsoleColor()
	}

	e := gin.New()

	// use logs middleware logurs
	// e.Use(gin.Logger())
	e.Use(logs.Logger(logs.GetLogger()))

	// use recovery middleware
	e.Use(gin.Recovery())

	// By default, http.ListenAndServe (which gin.Run wraps) will serve an unbounded number of requests.
	// Limiting the number of simultaneous connections can sometimes greatly speed things up under load
	if config.LimitConnections > 0 {
		e.Use(limit.MaxAllowed(config.LimitConnections))
	}

	return e
}

// Init engine init
func (config *EngineConfig) Init() http.Handler {
	e := config.initEngineConfig()
	// here init the routers, need refactor
	routers.InitRouters(e)
	return e
}
