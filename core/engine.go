/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package core

import (
	"net/http"
	"time"

	"github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/seeleteam/monitor-api/api/routers"
	"github.com/seeleteam/monitor-api/core/logs"
)

// EngineConfig gin engine config
type EngineConfig struct {
	DisableConsoleColor bool // disable the console color
	WriteLog            bool
	LogFile             string
	TempFolder          string
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

	// set gin mode release(hide handlers info)
	gin.SetMode(gin.ReleaseMode)

	e := gin.New()

	// use logs middleware logurs
	// e.Use(gin.Logger())
	e.Use(logs.Logger(logs.GetLogger()))

	// use recovery middleware
	e.Use(gin.Recovery())

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	corsConfig.AllowAllOrigins = true
	e.Use(cors.New(corsConfig))

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
