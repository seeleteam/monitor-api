/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package server

import (
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/seeleteam/monitor-api/config"
	"github.com/seeleteam/monitor-api/core"
	"github.com/seeleteam/monitor-api/core/logs"
	"github.com/seeleteam/monitor-api/rpc"
	"github.com/seeleteam/monitor-api/ws"
)

// Start WithErrorGroup
func Start(g *errgroup.Group) {
	// init the logger
	logs.NewLogger()

	monitorServer := core.GetServer(g)
	monitorServer.RunServer()

	// start RPCService, if enableWs = true
	enableWs := config.SeeleConfig.ServerConfig.EnableWebSocket
	if enableWs {
		logs.Infoln("will start web socket")
		time.Sleep(5 * time.Second)
		startWsService()
	} else {
		logs.Errorln("web socket start failed, EnableWebSocket is false")
		os.Exit(-1)
	}

}

func startWsService() {
	enableRPC := config.SeeleConfig.ServerConfig.EnableRPC
	if !enableRPC {
		logs.Fatalln("start RPC Service failed, EnableRPC is false")
		return
	}

	rpcURL := config.SeeleConfig.ServerConfig.RPCConfig.URL
	rpcSeeleRPC := rpc.NewSeeleRPC(rpcURL)

	wsURL := config.SeeleConfig.ServerConfig.WebSocketConfig.WsURL
	service, err := ws.New(wsURL, rpcSeeleRPC)
	if err != nil {
		log.Fatal(err)
	}
	//go service.Start()
	service.Start()
}
