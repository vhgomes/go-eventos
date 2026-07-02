package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vhgomes-eventos/internal/domain"
	erro "vhgomes-eventos/internal/pkg/errors"
	"vhgomes-eventos/internal/pkg/logger"
	"vhgomes-eventos/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"
)

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req LoginRequest) (string, error)
}

type RegisterRequest struct {
	Email    string
	Password string
	Name     string
}

type LoginRequest struct {
	Email    string
	Password string
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*domain.User, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("failed to check email existence in repository", err, zap.String("email", req.Email))
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		logger.Warn("registration failed, email already exists", zap.String("email", req.Email))
		return nil, erro.ErrConflict
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &domain.User{
		Email:    req.Email,
		Password: string(hashed),
		Name:     req.Name,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("failed to create user in repository", err, zap.String("email", req.Email))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, erro.ErrNotFound) {
			logger.Warn("login failed, user not found", zap.String("email", req.Email))
			return "", erro.ErrUnauthorized
		}
		logger.Error("failed to get user from repository", err, zap.String("email", req.Email))
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Warn("login failed, invalid password", zap.String("email", req.Email))
		return "", erro.ErrUnauthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.Id,
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}
