package http

import (
	"carWash/docs"
	"carWash/internal/config"
	v1 "carWash/internal/delivery/http/v1"
	"carWash/internal/service"
	"carWash/pkg/auth"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
	signingKey   string
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager, signingKey string) *Handler {
	return &Handler{services: services, tokenManager: tokenManager, signingKey: signingKey}
}

func (h *Handler) Init(cfg *config.Config) *fiber.App {
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	if cfg.Environment == config.Prod {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s", cfg.HTTP.Host)
	}
	router := fiber.New()
	router.Use(logger.New())
	router.Get("/swagger/*", swagger.HandlerDefault)

	h.initApi(router)
	router.Static("/media", "media")
	return router
}

func (h *Handler) initApi(router *fiber.App) {
	handler := v1.NewHandler(h.services, h.tokenManager, h.signingKey)
	api := router.Group("/api")
	{
		handler.Init(api)
	}
}
