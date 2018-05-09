/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package logs

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// New returns a gin compatable middleware using logrus to logs
// skipPaths only skips the INFO loglevel
func New(logger *logrus.Logger, skipPaths ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(skipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range skipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		c.Next()

		statusCode := c.Writer.Status()
		latency := time.Now().Sub(start)

		entry := logger.WithFields(logrus.Fields{
			"status":         statusCode,
			"method":         c.Request.Method,
			"path":           path,
			"ip":             c.ClientIP(),
			"latency":        latency,
			"latency_string": latency.String(),
			"user-agent":     c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
			return
		}

		if statusCode > 499 {
			entry.Error()
		} else if statusCode > 399 {
			entry.Warn()
		} else {
			if _, ok := skip[path]; ok {
				return
			}
			entry.Info()
		}

	}
}
