package service

import (
	"carWash/internal/domain"
	repos "carWash/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type FootServiceService struct {
	repos repos.FootService
}

func (f *FootServiceService) Update(ctx *fiber.Ctx, id int, inp domain.FootService) error {
	return f.repos.Update(ctx, id, inp)
}

func (f *FootServiceService) Create(ctx *fiber.Ctx, service domain.FootService) (int, error) {
	return f.repos.Create(ctx, service)
}

func (f *FootServiceService) GetAll(ctx *fiber.Ctx, page domain.Pagination) (*domain.GetAllResponses, error) {
	return f.repos.GetAll(ctx, page)
}

func (f *FootServiceService) GetById(ctx *fiber.Ctx, userId int) (*domain.FootService, error) {
	return f.repos.GetById(ctx, userId)
}

func (f *FootServiceService) Delete(ctx *fiber.Ctx, userId int) error {
	return f.repos.Delete(ctx, userId)
}

func NewFootServiceService(repos repos.FootService) *FootServiceService {
	return &FootServiceService{repos: repos}
}
