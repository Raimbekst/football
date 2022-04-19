package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
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

func (b *BuildingService) GetById(ctx *fiber.Ctx, id int) (*domain.Building, error) {
	return b.repos.GetById(ctx, id)
}

func (b *BuildingService) Update(ctx *fiber.Ctx, id int, inp domain.Building) error {
	return b.repos.Update(ctx, id, inp)
}

func (b *BuildingService) Delete(ctx *fiber.Ctx, id int) error {
	return b.repos.Delete(ctx, id)
}
