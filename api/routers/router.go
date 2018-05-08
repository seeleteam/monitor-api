/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/seeleteam/monitor-api/config"
)

// InitRouters init routers
func InitRouters(e *gin.Engine) {
	//web socket
	enableWs := config.SeeleConfig.ServerConfig.EnableWebSocket
	if enableWs {
		InitWsRouters(e)
	}
}
