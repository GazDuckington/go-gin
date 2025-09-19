package routes

import (
	"net/http"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/pkgs/auth"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, cfg *config.Config) {
	authGroup := r.Group("/auth")

	authGroup.POST("/login", func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			dto.Error(c, http.StatusBadRequest, "Invalid request payload", err.Error())
			return
		}

		// Dummy check for now
		if req.Username != "admin" || req.Password != "secret" {
			dto.Error(c, http.StatusUnauthorized, "Invalid username or password", nil)
			return
		}

		accessToken, err := auth.GenerateAccessToken(req.Username, "admin", cfg)
		if err != nil {
			dto.Error(c, http.StatusInternalServerError, "Failed to generate access token", err.Error())
			return
		}

		refreshToken, err := auth.GenerateRefreshToken(req.Username, cfg)
		if err != nil {
			dto.Error(c, http.StatusInternalServerError, "Failed to generate refresh token", err.Error())
			return
		}

		dto.Success(c, http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}, "Login successful")
	})
}
