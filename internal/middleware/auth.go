package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthConfig struct {
	Logger *logrus.Logger
	// Optionally add a function for token validation
	ValidateToken func(token string) (any, error)
}

// AuthRequired returns a middleware that enforces authentication.
// It extracts Bearer tokens from the Authorization header and validates them.
func AuthRequired(cfg AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			cfg.Logger.Warnf("unauthorized access to %s", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid authorization header",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		if cfg.ValidateToken == nil {
			cfg.Logger.Warn("no token validator configured, skipping validation")
			// You might want to reject here in production
			c.Next()
			return
		}

		claims, err := cfg.ValidateToken(token)
		if err != nil {
			cfg.Logger.WithError(err).Warn("token validation failed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Store claims in context for downstream handlers
		c.Set("authClaims", claims)
		c.Next()
	}
}
