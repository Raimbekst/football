package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type CardService struct {
	repos repository.Card
}

func NewCardService(repos repository.Card) *CardService {
	return &CardService{repos: repos}
}

func (f *CardService) Update(ctx *fiber.Ctx, id int, inp domain.Card) error {
	return f.repos.Update(ctx, id, inp)
}

func (f *CardService) Create(ctx *fiber.Ctx, card domain.Card) (int, error) {
	return f.repos.Create(ctx, card)
}

func (f *CardService) GetAll(ctx *fiber.Ctx, page domain.Pagination, userId int) (*domain.GetAllResponses, error) {
	return f.repos.GetAll(ctx, page, userId)
}

func (f *CardService) GetById(ctx *fiber.Ctx, id, userId int) (*domain.Card, error) {
	return f.repos.GetById(ctx, id, userId)
}

func (f *CardService) Delete(ctx *fiber.Ctx, id, userId int) error {
	return f.repos.Delete(ctx, id, userId)
}
