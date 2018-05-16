/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package filters

import (
	"github.com/gin-gonic/gin"

	"github.com/seeleteam/monitor-api/core/logs"
)

// BaseFilter base filter middleware
func BaseFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.ContentType()
		logs.Debug("contentType: %v", contentType)
		if contentType != gin.MIMEJSON {
			logs.Error("contentType must be %v", gin.MIMEJSON)
			c.Abort()
		}
		c.Next()
	}
}
