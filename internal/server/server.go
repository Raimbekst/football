package server

import (
	"carWash/internal/config"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	httpServer *fiber.App
}

func NewServer(handler *fiber.App) *Server {
	return &Server{
		httpServer: handler,
	}
}

func (s *Server) Run(cfg *config.Config) error {
	return fmt.Errorf("server.Run: %w", s.httpServer.Listen(fmt.Sprintf(":%s", cfg.HTTP.Port)))
}

func (s *Server) Stop() error {
	return fmt.Errorf("server.Stop: %w", s.httpServer.Shutdown())
}
