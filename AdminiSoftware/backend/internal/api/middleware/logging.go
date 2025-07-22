
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"status_code":  param.StatusCode,
			"latency":      param.Latency,
			"client_ip":    param.ClientIP,
			"method":       param.Method,
			"path":         param.Path,
			"user_agent":   param.Request.UserAgent(),
			"error":        param.ErrorMessage,
			"timestamp":    param.TimeStamp.Format(time.RFC3339),
		}).Info("HTTP Request")

		return ""
	})
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(400, gin.H{"error": "Invalid request format", "details": err.Error()})
			case gin.ErrorTypePublic:
				c.JSON(500, gin.H{"error": err.Error()})
			default:
				c.JSON(500, gin.H{"error": "Internal server error"})
			}
		}
	}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

func generateRequestID() string {
	// Simple request ID generation
	return time.Now().Format("20060102150405") + "-" + "abcdef"[:6]
}
package middleware

import (
	"AdminiSoftware/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(logger *utils.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Infof("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
		return ""
	})
}

func RequestLogger(logger *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := utils.GetClientIP(c.Request.RemoteAddr, c.GetHeader("X-Forwarded-For"), c.GetHeader("X-Real-IP"))
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Infof("[%d] %s %s - %s (%v)",
			statusCode,
			method,
			path,
			clientIP,
			latency,
		)
	}
}

func ErrorLogger(logger *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Errorf("Request error: %s", err.Error())
			}
		}
	}
}
