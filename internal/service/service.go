package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"carWash/pkg/auth"
	"carWash/pkg/hash"
	"carWash/pkg/phone"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"time"
)

type SignUpInput struct {
	Name            string
	PhoneNumber     string
	Password        string
	ConfirmPassword string
	UserType        string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type FavouriteInput struct {
	UserId     int
	BuildingId int
}

type UserAuth interface {
	VerifyExistenceUser(phone string) error
	UserSignUp(input SignUpInput) (string, error)
	Verify(input domain.VerifyUserInput) error

	UserSignIn(user domain.User) (*Tokens, error)

	SetPassword(id int, input domain.SetPasswordInput) error

	ResetPassword(phone string) (string, error)
	VerifyPhoneNumber(input domain.VerifyPhoneNumberInput) error
	ResetPasswordConfirm(input domain.ResetPasswordInput) error

	UpdatePhoneNumberVerify(inp domain.User) (string, error)
	UpdatePhoneNumberConfirm(input domain.ResetPhoneNumberInput, id int) error
}

type Building interface {
	Create(c *fiber.Ctx, building domain.Building) (int, error)
	GetAll(c *fiber.Ctx, page domain.Pagination, info domain.UserInfo, building domain.FilterForBuilding) (*domain.GetAllResponses, error)
	GetById(c *fiber.Ctx, id int) (*domain.Building, error)
	Update(c *fiber.Ctx, id int, inp domain.Building) error
	Delete(c *fiber.Ctx, id int) error
}

type BuildingImage interface {
	Create(c *fiber.Ctx, building domain.BuildingImage) (int, error)
	GetAll(c *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error)
	GetById(c *fiber.Ctx, id int) (*domain.BuildingImage, error)
	Update(c *fiber.Ctx, id int, inp domain.BuildingImage) error
	Delete(c *fiber.Ctx, id int) error
}

type Pitch interface {
	Create(ctx *fiber.Ctx, pitch domain.Pitch) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error)
	GetById(ctx *fiber.Ctx, id int) (*domain.Pitch, error)
	Update(ctx *fiber.Ctx, id int, inp domain.Pitch) error
	Delete(ctx *fiber.Ctx, id int) error
}

type Favourite interface {
	Create(ctx *fiber.Ctx, input FavouriteInput) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error)
	GetById(ctx *fiber.Ctx, id, userId int) (*domain.Favourite, error)
	Delete(ctx *fiber.Ctx, id, userId int) error
}

type Order interface {
	Create(ctx *fiber.Ctx, order domain.Order) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination, info domain.UserInfo, date float64) (*domain.GetAllResponses, error)
}

type Service struct {
	UserAuth
	Building
	BuildingImage
	Pitch
	Favourite
	Order
}

type Deps struct {
	Repos           *repository.Repository
	Hashes          hash.PasswordHashes
	OtpPhone        phone.SecretGenerator
	Ctx             context.Context
	Redis           *redis.Client
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewService(deps Deps) *Service {
	return &Service{
		UserAuth:      NewUserAuthService(deps.Repos.UserAuth, deps.Hashes, deps.OtpPhone, deps.Redis, deps.Ctx, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Building:      NewBuildingService(deps.Repos.Building),
		BuildingImage: NewBuildingImageService(deps.Repos.BuildingImage),
		Pitch:         NewPitchService(deps.Repos.Pitch),
		Favourite:     NewFavouriteService(deps.Repos.Favourite),
		Order:         NewOrderService(deps.Repos.Order),
	}
}
