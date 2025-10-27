package middleware

import (
	"fmt"
	"github.com/Asylann/Auth/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] %s %s %d %s \n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC822),
			params.Method,
			params.Path,
			params.StatusCode,
			params.Latency,
		)
	})
}

func RecoveryMiddle(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			// we are running inside the same goroutine where the panic happened
			log.Printf("recovered panic: %v\n%s", r, debug.Stack())
			c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
		}
	}()
	// call next middleware / handler in chain
	c.Next()
}

func RequireRole(role string, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIn, ok := c.Get("Role")
		logger.Infof("%s", roleIn)
		if !ok {
			logger.Error("Forbidden endpoint! Cant find role")
			c.IndentedJSON(http.StatusForbidden, response.Response{Err: "Forbidden endpoint!"})
			c.Abort()
			return
		}
		roleStr, ok := roleIn.(string)
		if !ok {
			c.IndentedJSON(http.StatusInternalServerError, response.Response{Err: "smt went wrong"})
			c.Abort()
			return
		}
		if roleStr != role {
			c.IndentedJSON(http.StatusForbidden, response.Response{Err: "Forbidden endpoint!"})
			c.Abort()
			return
		}
		c.Next()
	}
}
