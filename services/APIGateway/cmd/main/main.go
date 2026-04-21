package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sheelestun/WatchHRs-/internal/config"
	"github.com/sheelestun/WatchHRs-/internal/database"
	"github.com/sheelestun/WatchHRs-/internal/redis"
	"github.com/sheelestun/WatchHRs-/internal/repository/pg"
	redis2 "github.com/sheelestun/WatchHRs-/internal/repository/redis"
	"github.com/sheelestun/WatchHRs-/internal/service"
	"github.com/sheelestun/WatchHRs-/internal/web/handler"
	"github.com/sheelestun/WatchHRs-/internal/web/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Настройка логгера
	initLogger(&cfg.Logger)

	// Подключение к базе данных
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Подключение к Redis
	redisClient, err := redis.Connect(&cfg.Redis)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to Redis")
	}
	defer redisClient.Close()

	pgRepository := pg.NewPgRepository(db)
	redisCache := redis2.NewRedisCache(redisClient)

	// Инициализация валидатора
	validate := validator.New()

	//Инициализация сервисов
	authService := service.NewAuthService(pgRepository, redisCache)
	employeeService := service.NewEmployeeService(pgRepository, validate)
	imageService := service.NewImageService(pgRepository, validate)
	managerService := service.NewManagerService(pgRepository, validate)
	screenshotStatsService := service.NewScreenshotStatisticService(pgRepository, validate)
	workSessionService := service.NewWorkSessionService(pgRepository, validate)

	// Инициализация хендлеров
	authHandler := handler.NewAuthHandler(authService, cfg.SecretKey)
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	imageHandler := handler.NewImageHandler(imageService)
	managerHandler := handler.NewManagerHandler(managerService)
	screenshotStatsHandler := handler.NewScreenshotStatisticHandler(screenshotStatsService)
	workSessionHandler := handler.NewWorkSessionHandler(workSessionService)

	r := router.NewRouter()

	// Создание HTTP сервера
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		log.WithField("address", server.Addr).Info("HTTP server starting")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Fatal("HTTP server failed")
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server forced to shutdown")
	}

	log.Info("Server exited")
}

func initLogger(cfg *config.LoggerConfig) {
	switch cfg.Level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}

	switch cfg.Format {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	}

	switch cfg.File {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	}
}
