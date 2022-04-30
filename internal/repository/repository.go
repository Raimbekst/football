package repository

import (
	"carWash/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"math"
)

const (
	userTable          = "users"
	sessionTable       = "sessions"
	buildingTable      = "buildings"
	buildingImageTable = "images"
	pitchTable         = "pitches"
	favouriteTable     = "favourites"
	orderTable         = "orders"
)

type FavouriteInput struct {
	UserId     int
	BuildingId int
}
type UserAuth interface {
	VerifyExistenceUser(phone string, activated bool) (*domain.User, error)

	UpdateUser(user domain.User, id int) error
	CreateUser(user domain.User) (int, error)
	Verify(phone string) error

	SignIn(phone, password string) (*domain.User, error)
	SetSession(userId int, session domain.Session) error

	GetUser(id int) (*domain.User, error)

	VerifyViaPassword(id int, password string) error
	SetPassword(id int, hashedOldPassword, hashedNewPassword string) error

	VerifyViaPhoneNumber(phone string) (*domain.User, error)
	ResetPassword(phone, password string) error
}

type Building interface {
	Create(c *fiber.Ctx, building domain.Building) (int, error)
	GetAll(c *fiber.Ctx, page domain.Pagination, info domain.UserInfo, building domain.FilterForBuilding) (*domain.GetAllResponses, error)
	GetById(c *fiber.Ctx, id int) (*domain.Building, error)
	Update(c *fiber.Ctx, id int, inp domain.Building) ([]string, error)
	Delete(c *fiber.Ctx, id int) ([]string, error)
}

type Pitch interface {
	Create(ctx *fiber.Ctx, pitch domain.Pitch) (int, error)
	GetAll(ctx *fiber.Ctx, page domain.Pagination, id int) (*domain.GetAllResponses, error)
	GetById(ctx *fiber.Ctx, id int) (*domain.Pitch, error)
	Update(ctx *fiber.Ctx, id int, inp domain.Pitch) ([]string, error)
	Delete(ctx *fiber.Ctx, id int) ([]string, error)
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

type Repository struct {
	UserAuth
	Building
	Pitch
	Favourite
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserAuth:  NewUserAuthRepos(db),
		Building:  NewBuildingRepos(db),
		Pitch:     NewPitchRepos(db),
		Favourite: NewFavouriteRepos(db),
		Order:     NewOrderRepos(db),
	}
}

func calculatePagination(page *domain.Pagination, count int) (int, int) {
	if page.Limit == 0 {
		page.Limit = count
	}

	if page.Page == 0 {
		page.Page = 1
	}

	pagesCount := 1.0

	if count != 0 {
		pagesCount = math.Ceil(float64(count) / float64(page.Limit))
		if page.Limit >= count {
			pagesCount = 1
		}
	}

	offset := (page.Page - 1) * page.Limit

	return offset, int(pagesCount)
}
