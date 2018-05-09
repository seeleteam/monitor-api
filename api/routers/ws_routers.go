/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package routers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/seeleteam/monitor-api/core/logs"
	"github.com/seeleteam/monitor-api/core/utils"
)

// InitWsRouters init the web socket api
func InitWsRouters(e *gin.Engine) {
	e.GET("/api", bindWsHandler())
}

// bindWsHandler bind the handler for web socket
func bindWsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	}
}

// web socket default config
var upGrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin:      func(r *http.Request) bool { return true },
	HandshakeTimeout: time.Duration(time.Second * 60),
}

// wsHandler web socket handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		logs.Error("Failed to set websocket upgrade: %+v", err)
		return
	}
	for {
		msgType, msgData, err := conn.ReadMessage()
		switch err.(type) {
		case *websocket.CloseError:
			logs.Error("websocket is closed, error is %v", err)
			return
		default:
		}
		if err != nil {
			break
		}

		// Skip binary messages
		if msgType != websocket.TextMessage {
			continue
		}

		var msg map[string][]interface{}
		err = json.Unmarshal(msgData, &msg)
		if err != nil {
			logs.Error("WsHandler receive msg %v", string(msgData))
		}

		resultData := make(map[string][]interface{})

		command, ok := msg["emit"][0].(string)
		if !ok {
			logs.Error("Invalid stats server message type", "type", msg["emit"][0])
			return
		}
		logs.Debug("receive msg len is %v, msg is %+v\n", len(msg["emit"]), utils.StructSerialize(msg))
		if len(msg["emit"]) == 2 && command == "node-ping" {
			hostname, _ := os.Hostname()
			resultData = map[string][]interface{}{
				"emit": {"node-pong", map[string]string{
					"id":         hostname + "_" + conn.LocalAddr().String(),
					"clientTime": time.Now().String(),
				}},
			}
			// write
			responseData := utils.StructSerialize(resultData)
			conn.WriteMessage(msgType, []byte(responseData))
			logs.Debug("output message, type(1=text, 2=binary): %+v, msg: %+v\n", msgType, responseData)
		} else {
			// write
			conn.WriteMessage(msgType, msgData)
			logs.Debug("output message, type(1=text, 2=binary): %+v, msg: %+v\n", msgType, string(msgData))
		}
	}
}
