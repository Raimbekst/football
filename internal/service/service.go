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

	GetUserInfo(id int) (*domain.User, error)
	UpdateUserInfo(user domain.UserUpdate, id int) error
}

type Building interface {
	Create(c *fiber.Ctx, building domain.Building) (int, error)
	GetAll(c *fiber.Ctx, page domain.Pagination, info domain.UserInfo, building domain.FilterForBuilding) (*domain.GetAllResponses, error)
	GetById(c *fiber.Ctx, info domain.UserInfo, id int) (*domain.Building, error)
	Update(c *fiber.Ctx, id int, inp domain.Building) error
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
	GetById(ctx *fiber.Ctx, id, userId int) (*domain.Building, error)
	Delete(ctx *fiber.Ctx, id, userId int) error
}

type Order interface {
	Create(ctx *fiber.Ctx, order domain.Order) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination, info domain.UserInfo, order domain.FilterForOrder) (*domain.GetAllResponses, error)
	GetAllBookTime(ctx *fiber.Ctx, times domain.FilterForOrderTimes) (*domain.GetAllResponses, error)
}

type Comment interface {
	Create(ctx *fiber.Ctx, comment domain.Comment) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination, buildingId int) (*domain.GetAllResponses, error)
	CreateGrade(ctx *fiber.Ctx, grade domain.Grade) (int, error)
	GetAllGrades(ctx *fiber.Ctx, page domain.Pagination, buildingId int) (*domain.GetAllResponses, error)
}

type Feedback interface {
	Create(ctx *fiber.Ctx, feedback domain.Feedback, id int) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination) (*domain.GetAllResponses, error)
}

type FootService interface {
	Create(ctx *fiber.Ctx, service domain.FootService) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination) (*domain.GetAllResponses, error)
	GetById(ctx *fiber.Ctx, userId int) (*domain.FootService, error)
	Update(ctx *fiber.Ctx, id int, inp domain.FootService) error
	Delete(ctx *fiber.Ctx, userId int) error
}

type Card interface {
	Create(ctx *fiber.Ctx, service domain.Card) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination, userId int) (*domain.GetAllResponses, error)
	GetById(ctx *fiber.Ctx, id, userId int) (*domain.Card, error)
	Update(ctx *fiber.Ctx, id int, inp domain.Card) error
	Delete(ctx *fiber.Ctx, id, userId int) error
}

type Service struct {
	UserAuth
	Building
	Pitch
	Favourite
	Order
	Comment
	Feedback
	FootService
	Card
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
		UserAuth:    NewUserAuthService(deps.Repos.UserAuth, deps.Hashes, deps.OtpPhone, deps.Redis, deps.Ctx, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Building:    NewBuildingService(deps.Repos.Building),
		Pitch:       NewPitchService(deps.Repos.Pitch),
		Favourite:   NewFavouriteService(deps.Repos.Favourite),
		Order:       NewOrderService(deps.Repos.Order),
		Comment:     NewCommentService(deps.Repos.Comment),
		Feedback:    NewFeedbackService(deps.Repos.Feedback),
		FootService: NewFootServiceService(deps.Repos.FootService),
		Card:        NewCardService(deps.Repos.Card),
	}
}
