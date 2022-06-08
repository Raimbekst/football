package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type OrderService struct {
	repos repository.Order
}

func (o *OrderService) GetAllBookTime(ctx *fiber.Ctx, times domain.FilterForOrderTimes) (*domain.GetAllResponses, error) {
	return o.repos.GetAllBookTime(ctx, times)
}

func NewOrderService(repos repository.Order) *OrderService {
	return &OrderService{repos: repos}
}

func (o *OrderService) Create(ctx *fiber.Ctx, order domain.Order) (int, error) {
	return o.repos.Create(ctx, order)
}

func (o *OrderService) GetAll(ctx *fiber.Ctx, page domain.Pagination, info domain.UserInfo, order domain.FilterForOrder) (*domain.GetAllResponses, error) {
	return o.repos.GetAll(ctx, page, info, order)
}
