package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"carWash/pkg/media"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type BuildingService struct {
	repos repository.Building
}

func NewBuildingService(repos repository.Building) *BuildingService {
	return &BuildingService{repos: repos}
}

func (b *BuildingService) Create(ctx *fiber.Ctx, building domain.Building) (int, error) {
	return b.repos.Create(ctx, building)
}
func (b *BuildingService) GetAll(ctx *fiber.Ctx, page domain.Pagination, info domain.UserInfo, building domain.FilterForBuilding) (*domain.GetAllResponses, error) {
	return b.repos.GetAll(ctx, page, info, building)
}

func (b *BuildingService) GetById(ctx *fiber.Ctx, info domain.UserInfo, id int) (*domain.Building, error) {
	return b.repos.GetById(ctx, info, id)
}

func (b *BuildingService) Update(ctx *fiber.Ctx, id int, inp domain.Building) error {
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

func (b *BuildingService) Delete(ctx *fiber.Ctx, id int) error {
	img, err := b.repos.Delete(ctx, id)
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
