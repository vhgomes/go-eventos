package handlers

import (
	"net/http"
	"vhgomes-eventos/internal/pkg/logger"
	"vhgomes-eventos/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body registerRequest true "User data"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("failed to bind register request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("starting user registration", zap.String("email", req.Email))

	user, err := h.authService.Register(c.Request.Context(), service.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		logger.Error("failed to register user", err, zap.String("email", req.Email))
		ResponseError(c, err)
		return
	}

	logger.Info("user registered successfully", zap.Int("user_id", user.Id))
	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body loginRequest true "Credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("failed to bind login request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Info("starting user login", zap.String("email", req.Email))

	token, err := h.authService.Login(c.Request.Context(), service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		logger.Error("login failed", err, zap.String("email", req.Email))
		ResponseError(c, err)
		return
	}

	logger.Info("user logged in successfully", zap.String("email", req.Email))
	c.JSON(http.StatusOK, loginResponse{Token: token})
}
