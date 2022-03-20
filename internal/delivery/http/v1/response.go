package v1

import (
	"carWash/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

const (
	user    = "user"
	admin   = "admin"
	manager = "manager"
)

type idResponse struct {
	ID interface{} `json:"id"`
}
type okResponse struct {
	Message string `json:"message"`
}

type response struct {
	Message string `json:"detail"`
}

func newResponse(c *fiber.Ctx, statusCode, message string) {
	logger.Error(message)

}
