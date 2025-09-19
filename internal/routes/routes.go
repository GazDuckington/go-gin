package routes

import (
	"time"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger(cfg))

	// health
	healthCtrl := controller.NewHealthController(cfg)
	r.GET("/health", healthCtrl.HealthCheck)

	// NOTE: register domains
	// RegisterUserRoutes(r, cfg)
	return r
}

func requestLogger(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		status := c.Writer.Status()

		entry := cfg.Logger.WithFields(map[string]interface{}{
			"status":     status,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency_ms": latency,
		})

		if len(c.Errors) > 0 {
			entry.WithField("errors", c.Errors.String()).Error("request error")
			return
		}

		switch {
		case status >= 500:
			entry.Error("internal error")
		case status >= 400:
			entry.Warn("client error")
		default:
			entry.Info("ok")
		}
	}
}
