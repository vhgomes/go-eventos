package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "vhgomes-eventos/docs"
	"vhgomes-eventos/internal/handlers"
	"vhgomes-eventos/internal/middleware"
	"vhgomes-eventos/internal/pkg/config"
	"vhgomes-eventos/internal/repository/postgres"
	"vhgomes-eventos/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3" // ou "github.com/lib/pq" para PostgreSQL
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("database unreachable: %v", err)
	}
	log.Printf("postgres connection established (max_open=%d, max_idle=%d)", cfg.MaxOpenConns, cfg.MaxIdleConns)

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
		log.Printf("starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("server error: %v", err)
	case sig := <-shutdown:
		log.Printf("shutting down server (signal: %s)", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
		if err := srv.Close(); err != nil {
			log.Printf("failed to close server: %v", err)
		}
	} else {
		log.Println("server stopped gracefully")
	}
}
