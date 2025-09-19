package controller

import (
	"net/http"
	"strconv"

	"github.com/GazDuckington/go-gin/internal/config"
	"github.com/GazDuckington/go-gin/internal/models/dto"
	"github.com/GazDuckington/go-gin/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	svc service.UserService
	cfg *config.Config
}

func NewUserController(s service.UserService, cfg *config.Config) *UserController {
	return &UserController{svc: s, cfg: cfg}
}

func (ctrl *UserController) GetAll(c *gin.Context) {
	ctrl.cfg.Logger.Debug("GetAll users called")
	users, err := ctrl.svc.GetAll()
	if err != nil {
		ctrl.cfg.Logger.Errorf("GetAll error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	user, err := ctrl.svc.GetByID(uint(id64))
	if err != nil {
		ctrl.cfg.Logger.Errorf("GetByID error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.cfg.Logger.Warnf("Create bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := ctrl.svc.Create(req)
	if err != nil {
		ctrl.cfg.Logger.Errorf("Create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusCreated, created)
}
