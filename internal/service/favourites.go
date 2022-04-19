package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type FavouriteService struct {
	repos repository.Favourite
}

func NewFavouriteService(repos repository.Favourite) *FavouriteService {
	return &FavouriteService{repos: repos}
}

func (f *FavouriteService) Create(ctx *fiber.Ctx, input FavouriteInput) (int, error) {
	return f.repos.Create(ctx, repository.FavouriteInput{
		UserId:     input.UserId,
		BuildingId: input.BuildingId,
	})
}

func (f *FavouriteService) GetAll(ctx *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error) {
	return f.repos.GetAll(ctx, page, id)
}

func (f *FavouriteService) GetById(ctx *fiber.Ctx, id, userId int) (*domain.Favourite, error) {
	return f.repos.GetById(ctx, id, userId)
}

func (f *FavouriteService) Delete(ctx *fiber.Ctx, id, userId int) error {
	return f.repos.Delete(ctx, id, userId)
}
