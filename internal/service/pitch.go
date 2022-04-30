package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"carWash/pkg/media"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type PitchService struct {
	repos repository.Pitch
}

func NewPitchService(repos repository.Pitch) *PitchService {
	return &PitchService{repos: repos}
}

func (p *PitchService) Create(ctx *fiber.Ctx, pitch domain.Pitch) (int, error) {
	return p.repos.Create(ctx, pitch)
}

func (p *PitchService) GetAll(ctx *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error) {
	return p.repos.GetAll(ctx, page, id)
}

func (p *PitchService) GetById(ctx *fiber.Ctx, id int) (*domain.Pitch, error) {
	return p.repos.GetById(ctx, id)
}

func (p *PitchService) Update(ctx *fiber.Ctx, id int, inp domain.Pitch) error {
	img, err := p.repos.Update(ctx, id, inp)
	if err != nil {
		return fmt.Errorf("service.Update: %w", err)

	}
	for i, _ := range img {
		if img[i] != "" {
			err = media.DeleteImage(img[i])
			if err != nil {
				return fmt.Errorf("service.Update: %w", err)
			}
		}
	}
	return nil
}

func (p *PitchService) Delete(ctx *fiber.Ctx, id int) error {
	img, err := p.repos.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("service.Delete: %w", err)

	}
	for i, _ := range img {
		if img[i] != "" {
			err = media.DeleteImage(img[i])
			if err != nil {
				return fmt.Errorf("service.Delete: %w", err)
			}
		}
	}
	return nil
}
