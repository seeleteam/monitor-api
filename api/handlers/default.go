/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package handlers

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seeleteam/monitor-api/core/logs"
)

// Ping defaut for test
func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, H{
			"message": "ping" + c.Request.URL.Path,
		})
	}
}

// Pong defaut for test
func Pong() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(204, H{
			"message": "pong" + c.Request.URL.Path,
		})
	}
}

// Kong defaut for test
func Kong() gin.HandlerFunc {
	return func(c *gin.Context) {
		code, err := c.Writer.WriteString("Kong" + c.Request.URL.Path)
		if err != nil {
			logs.Error("error is %v %v", code, err)
			logs.Errorln("error is ", code, err)

		} else {
			logs.Info("info is %v and time now %v", code, time.Now())
			logs.Infoln("info is ", code, time.Now().Nanosecond())
		}
	}
}

// LongAsync async task with goroutine!
func LongAsync() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		go func() {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)

			// note that you are using the copied context "cCp", IMPORTANT
			logs.Printf("Done! in path %v", cCp.Request.URL.Path)
		}()
	}
}
