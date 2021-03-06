package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type FeedbackService struct {
	repos repository.Feedback
}

func (f *FeedbackService) CreateNoty(ctx *fiber.Ctx, feedback domain.Notification) (int, error) {
	return f.repos.CreateNoty(ctx, feedback)
}

func (f *FeedbackService) GetAllNoty(ctx *fiber.Ctx, page domain.Pagination) (*domain.GetAllResponses, error) {
	return f.repos.GetAllNoty(ctx, page)
}

func NewFeedbackService(repos repository.Feedback) *FeedbackService {
	return &FeedbackService{repos: repos}
}

func (f *FeedbackService) Create(ctx *fiber.Ctx, feedback domain.Feedback, id int) (int, error) {
	return f.repos.Create(ctx, feedback, id)
}

func (f *FeedbackService) GetAll(ctx *fiber.Ctx, page domain.Pagination) (*domain.GetAllResponses, error) {
	return f.repos.GetAll(ctx, page)
}
