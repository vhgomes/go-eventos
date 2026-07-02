package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "vhgomes-eventos/docs"
	"vhgomes-eventos/internal/handlers"
	"vhgomes-eventos/internal/middleware"
	"vhgomes-eventos/internal/pkg/config"
	"vhgomes-eventos/internal/pkg/logger"
	"vhgomes-eventos/internal/repository/postgres"
	"vhgomes-eventos/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to database", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verifica conectividade
	if err := db.Ping(); err != nil {
		logger.Fatal("database unreachable", err)
	}
	logger.Info("postgres connection established", zap.Int("max_open", cfg.MaxOpenConns), zap.Int("max_idle", cfg.MaxIdleConns))

	userRepo := postgres.NewUserRepo(db)
	eventRepo := postgres.NewEventRepo(db)
	attendeeRepo := postgres.NewAttendeeRepo(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	eventService := service.NewEventService(eventRepo)
	attendeeService := service.NewAttendeeService(attendeeRepo, eventRepo, userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	eventHandler := handlers.NewEventHandler(eventService)
	attendeeHandler := handlers.NewAttendeeHandler(attendeeService)

	router := gin.Default()

	router.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.URL.Path == "/swagger/" {
			c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/events", eventHandler.GetAll)
		v1.GET("/events/:id", eventHandler.GetByID)
		v1.GET("/events/:id/attendees", attendeeHandler.GetAttendees)
		v1.GET("/attendees/:id/events", attendeeHandler.GetEventsByAttendee)
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
	}

	authGroup := v1.Group("/")
	authGroup.Use(middleware.AuthMiddleware(userRepo, cfg.JWTSecret))
	{
		authGroup.POST("/events", eventHandler.Create)
		authGroup.PUT("/events/:id", eventHandler.Update)
		authGroup.DELETE("/events/:id", eventHandler.Delete)
		authGroup.POST("/events/:id/attendees/:userId", attendeeHandler.AddAttendee)
		authGroup.DELETE("/events/:id/attendees/:userId", attendeeHandler.RemoveAttendee)
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Fatal("server error", err)
	case sig := <-shutdown:
		logger.Info("shutting down server", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", err)
		if err := srv.Close(); err != nil {
			logger.Error("failed to close server", err)
		}
	} else {
		logger.Info("server stopped gracefully")
	}
}
