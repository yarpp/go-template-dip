package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"advanced-blog-management-system/internal/handler"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/service"
	"advanced-blog-management-system/internal/storage"
	"advanced-blog-management-system/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// .env (не критично, если нет)
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env not found, using environment variables")
	}

	// Конфигурация
	serverHost := getEnv("SERVER_HOST", "0.0.0.0")
	serverPort := getEnv("SERVER_PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "super-secret-key")
	jwtExpiryStr := getEnv("JWT_EXPIRY_HOURS", "24")
	dataDir := getEnv("DATA_DIR", "./data")
	logFilePath := getEnv("LOGS_FILE", "./log.txt")

	jwtExpiry, _ := strconv.Atoi(jwtExpiryStr)
	if jwtExpiry <= 0 {
		jwtExpiry = 24
	}

	// Хранилища
	userStore, err := storage.NewUserStore(dataDir + "/users.json")
	if err != nil {
		log.Fatalf("failed to init user store: %v", err)
	}
	postStore, err := storage.NewPostStore(dataDir + "/posts.json")
	if err != nil {
		log.Fatalf("failed to init post store: %v", err)
	}
	commentStore, err := storage.NewCommentStore(dataDir + "/comments.json")
	if err != nil {
		log.Fatalf("failed to init comment store: %v", err)
	}

	// Event logger
	eventLogger := logger.NewEventLogger(logFilePath)
	eventLogger.Start()

	// JWT
	jwtManager := auth.NewJWTManager(jwtSecret, jwtExpiry)

	// Сервисы
	userService := service.NewUserService(userStore, jwtManager)
	postService := service.NewPostService(postStore, userStore, eventLogger)
	commentService := service.NewCommentService(commentStore, postStore, eventLogger)

	// Handlers
	authHandler := handler.NewAuthHandler(userService)
	postHandler := handler.NewPostHandler(postService, eventLogger)
	commentHandler := handler.NewCommentHandler(commentService, eventLogger)

	// Middleware
	loggingMW := middleware.NewLoggingMiddleware(log.New(os.Stdout, "[http] ", log.LstdFlags))
	authMW := middleware.NewAuthMiddleware(jwtManager)

	router := setupRouter(authHandler, postHandler, commentHandler, loggingMW, authMW)

	srv := &http.Server{
		Addr:         serverHost + ":" + serverPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	// Останавливаем логгер и дожидаемся записи
	eventLogger.Stop()
	log.Println("server stopped")
}

func setupRouter(
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler,
	loggingMW *middleware.LoggingMiddleware,
	authMW *middleware.AuthMiddleware,
) *chi.Mux {
	r := chi.NewRouter()

	// Глобальные middleware
	r.Use(loggingMW.Recovery)
	r.Use(loggingMW.Logger)
	r.Use(loggingMW.CORS)

	r.Route("/api", func(api chi.Router) {
		api.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		})

		// Публичные
		api.Post("/register", authHandler.Register)
		api.Post("/login", authHandler.Login)

		api.Get("/posts", postHandler.GetAll)
		api.Get("/posts/{id}", postHandler.GetByID)
		api.Get("/posts/{postId}/comments", commentHandler.GetByPost)
		api.Get("/posts/author/{authorID}", postHandler.GetByAuthor)

		// Защищённые
		api.Group(func(priv chi.Router) {
			priv.Use(middleware.ToMiddleware(authMW.RequireAuth))
			priv.Post("/posts", postHandler.Create)
			priv.Post("/posts/{postId}/comments", commentHandler.Create)
		})
	})

	return r
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
