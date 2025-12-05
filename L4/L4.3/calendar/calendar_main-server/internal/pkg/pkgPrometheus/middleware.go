package pkgPrometheus

import (
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	reUUID   = regexp.MustCompile(`[0-9a-fA-F-]{36}`)
	reNumber = regexp.MustCompile(`/\d+`)
	reHexID  = regexp.MustCompile(`/[0-9a-f]{8,}`)
)

func normalizePath(path string) string {
	path = reUUID.ReplaceAllString(path, "/:uuid")
	path = reHexID.ReplaceAllString(path, "/:hex")
	path = reNumber.ReplaceAllString(path, "/:id")
	return path
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		statusCode := c.Writer.Status()
		method := c.Request.Method

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		path = normalizePath(path)

		RequestCount.WithLabelValues(
			strconv.Itoa(statusCode),
			method,
			path,
		).Inc()

		RequestDuration.WithLabelValues(
			method,
			path,
		).Observe(duration)
	}
}
