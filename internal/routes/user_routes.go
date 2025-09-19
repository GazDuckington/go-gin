package routes

import (
	database "github.com/GazDuckington/go-gin/db"
	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/controller"
	middleware "github.com/GazDuckington/go-gin/internal/middlewares"
	"github.com/GazDuckington/go-gin/internal/repository"
	"github.com/GazDuckington/go-gin/internal/service"
	"github.com/GazDuckington/go-gin/pkgs/auth"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, cfg *config.Config) {
	// wire concrete implementations
	userRepo := repository.NewUserRepository(database.DB)
	userSvc := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userSvc, cfg)

	g := r.Group("/users")
	g.Use(middleware.AuthRequired(middleware.AuthConfig{
		Logger:        cfg.Logger,
		ValidateToken: auth.ValidateJWT,
	}))
	{
		g.GET("", userCtrl.GetAll)
		g.POST("", userCtrl.Create)
		g.GET("/:id", userCtrl.GetByID)
	}
}
