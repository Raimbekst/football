package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type CommentService struct {
	repos repository.Comment
}

func NewCommentService(repos repository.Comment) *CommentService {
	return &CommentService{repos: repos}
}

func (c *CommentService) GetAll(ctx *fiber.Ctx, page domain.Pagination, buildingId int) (*domain.GetAllResponses, error) {
	return c.repos.GetAll(ctx, page, buildingId)
}

func (c *CommentService) Create(ctx *fiber.Ctx, comment domain.Comment) (int, error) {
	return c.repos.Create(ctx, comment)
}
