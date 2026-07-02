package middleware

import (
	"net/http"
	"strings"
	"vhgomes-eventos/internal/pkg/logger"
	"vhgomes-eventos/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func AuthMiddleware(userRepo repository.UserRepository, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("authorization header required", zap.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			logger.Warn("bearer token required", zap.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "bearer token required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			logger.Warn("invalid token", zap.Error(err), zap.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn("invalid token claims", zap.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["userId"].(float64)
		if !ok {
			logger.Warn("invalid token claims, missing userId", zap.String("client_ip", c.ClientIP()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		userID := int(userIDFloat)

		user, err := userRepo.GetByID(c.Request.Context(), userID)
		if err != nil {
			logger.Warn("user not found from token", zap.Int("user_id", userID), zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		logger.Debug("user authenticated successfully", zap.Int("user_id", user.Id))
		c.Set("user", user)
		c.Next()
	}
}
