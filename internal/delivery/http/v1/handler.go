package v1

import (
	"carWash/internal/service"
	"carWash/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

const (
	reserved = 1
	finished = 2
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
	signingKey   string
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager, signingKey string) *Handler {
	return &Handler{services: services, tokenManager: tokenManager, signingKey: signingKey}
}

func (h *Handler) Init(api fiber.Router) {
	v1 := api.Group("/v1")
	{
		h.initUserRoutes(v1)
		h.initOrderRoutes(v1)
		h.initBuildingRoutes(v1)
		h.initPitchRoutes(v1)
		h.initFavouriteRoutes(v1)
	}
}

func getUser(c *fiber.Ctx) (string, int) {

	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)

	id := claims["jti"].(string)

	userType := claims["sub"].(string)

	idInt, err := strconv.Atoi(id)

	if err != nil {
		return "", 0
	}
	return userType, idInt
}

func isAdmin(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	return c.Next()
}

func isUser(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "user" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	return c.Next()
}

func isManager(c *fiber.Ctx) error {
	userType, _ := getUser(c)

	if userType != "manager" {
		return c.Status(fiber.StatusUnauthorized).JSON(response{Message: "нет доступа"})
	}
	return c.Next()
}
