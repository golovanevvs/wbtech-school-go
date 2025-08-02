package handler

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type responseBody struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseBody) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (h Handler) WithLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		logLevel := h.Logger.GetLevel()

		if logLevel == zerolog.Disabled {
			c.Next()
			return
		}

		start := time.Now()

		var writer *responseBody
		if logLevel <= zerolog.DebugLevel {
			writer = &responseBody{
				ResponseWriter: c.Writer,
				body:           bytes.NewBufferString(""),
			}
			c.Writer = writer
		}

		c.Next()

		logCtx := h.Logger.With().
			Str("Request method", c.Request.Method).
			Str("Request path", c.Request.URL.Path).
			Str("Request query", c.Request.URL.RawQuery).
			Str("Request ip", c.ClientIP()).
			Str("Request Content-Type", c.ContentType()).
			Str("Request user-agent", c.Request.UserAgent()).
			Int("Response status", c.Writer.Status()).
			Int("Response size", c.Writer.Size())

		if logLevel <= zerolog.DebugLevel && writer != nil {
			if body := writer.body.String(); body != "" {
				if strings.Contains(writer.Header().Get("Content-Type"), "application/json") {
					var pretty bytes.Buffer
					if err := json.Indent(&pretty, []byte(body), "", "  "); err != nil {
						logCtx = logCtx.RawJSON("Response body", []byte(body))
					} else {
						logCtx = logCtx.RawJSON("Response body", []byte(pretty.Bytes()))
					}
				} else {
					var truncateBody string
					maxLenBody := 1024
					if len(body) > maxLenBody {
						truncateBody = body[:maxLenBody] + "...[truncated]"
					} else {
						truncateBody = body
					}
					logCtx = logCtx.Str("Response body", truncateBody)
				}
			}
		}

		logCtx = logCtx.Dur("latency", time.Since(start))

		log := logCtx.Logger()

		msg := "Request handled"

		if len(c.Errors) > 0 {
			errors := make([]error, len(c.Errors))
			for i, e := range c.Errors {
				errors[i] = e.Err
			}
			log.Error().Errs("errors", errors).Msg(msg)
		} else {
			switch {
			case logLevel <= zerolog.DebugLevel:
				log.Debug().Msg(msg)
			case logLevel <= zerolog.InfoLevel:
				log.Info().Msg(msg)
			}
		}
	}
}
