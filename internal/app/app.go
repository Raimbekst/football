package app

import (
	"carWash/internal/config"
	delivery "carWash/internal/delivery/http"
	repos "carWash/internal/repository"
	"carWash/internal/server"
	"carWash/internal/service"
	"carWash/pkg/auth"
	"carWash/pkg/database"
	"carWash/pkg/database/redis"
	"carWash/pkg/hash"
	"carWash/pkg/logger"
	"carWash/pkg/phone"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error(err)
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Error(err)
	}
	red, err := redis.NewRedisDB(cfg)
	if err != nil {
		logger.Error(err)
	}

	hashes := hash.NewSHA1Hashes(cfg.Auth.PasswordSalt)

	otpNumberGenerator := phone.NewSecretGenerator()

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
	}

	repository := repos.NewRepository(db)
	services := service.NewService(service.Deps{
		Repos:           repository,
		Hashes:          hashes,
		OtpPhone:        otpNumberGenerator,
		Ctx:             context.TODO(),
		Redis:           red,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.JWT.RefreshTokenTTL,
	})

	handlers := delivery.NewHandler(services, tokenManager, cfg.Auth.JWT.SigningKey)

	srv := server.NewServer(handlers.Init(cfg))

	go func() {
		if err := srv.Run(cfg); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occur while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("server started")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 3 * time.Second

	_, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(); err != nil {
		logger.Errorf("failed to stop server: %v", err)
	}
}
