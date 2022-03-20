package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"carWash/pkg/media"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type BuildingImageService struct {
	repos repository.BuildingImage
}

func NewBuildingImageService(repos repository.BuildingImage) *BuildingImageService {
	return &BuildingImageService{repos: repos}
}

func (b *BuildingImageService) Create(ctx *fiber.Ctx, building domain.BuildingImage) (int, error) {
	return b.repos.Create(ctx, building)
}

func (b *BuildingImageService) GetAll(ctx *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error) {
	return b.repos.GetAll(ctx, page, id)
}

func (b *BuildingImageService) GetById(ctx *fiber.Ctx, id int) (*domain.BuildingImage, error) {
	return b.repos.GetById(ctx, id)
}

func (b *BuildingImageService) Update(ctx *fiber.Ctx, id int, inp domain.BuildingImage) error {
	img, err := b.repos.Update(ctx, id, inp)
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

func (b *BuildingImageService) Delete(ctx *fiber.Ctx, id int) error {
	img, err := b.repos.Delete(ctx, id)
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
