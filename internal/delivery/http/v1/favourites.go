package v1

import (
	"carWash/internal/domain"
	"carWash/internal/service"
	"carWash/pkg/validation/validationStructs"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
)

type Favourite struct {
	BuildingId int `json:"building_id" validate:"required"`
}

func (h *Handler) initFavouriteRoutes(api fiber.Router) {
	router := api.Group("/favourite", jwtware.New(
		jwtware.Config{
			SigningKey: []byte(h.signingKey)}), isUser)
	{
		router.Post("", h.createFavourite)
		router.Get("", h.getAllFavourites)
		router.Get("/:id", h.getFavouriteById)
		router.Delete("/:id", h.deleteFavourite)
	}
}

// @Security User_Auth
// @Tags favourite
// @Description create favourite
// @Accept json
// @Produce json
// @Param data body Favourite true "favourite create input"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /favourite [post]
func (h *Handler) createFavourite(c *fiber.Ctx) error {
	var (
		inp Favourite
	)

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(inp)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}
	_, userId := getUser(c)

	id, err := h.services.Favourite.Create(c, service.FavouriteInput{
		UserId:     userId,
		BuildingId: inp.BuildingId,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{id})

}

// @Security User_Auth
// @Tags favourite
// @Description get all favourites
// @ID get-all-favourites
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /favourite [get]
func (h *Handler) getAllFavourites(c *fiber.Ctx) error {
	var (
		page domain.Pagination
	)

	_, userId := getUser(c)

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	list, err := h.services.Favourite.GetAll(c, page, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)

}

// @Security User_Auth
// @Tags favourite
// @Description get favourite by id
// @ID get-favourite-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "favourite id"
// @Success 200 {object} domain.Favourite
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /favourite/{id} [get]
func (h *Handler) getFavouriteById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, userId := getUser(c)

	list, err := h.services.Favourite.GetById(c, id, userId)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags favourite
// @Description delete favourite
// @ModuleID deleteFavourite
// @Accept  json
// @Produce  json
// @Param id path string true "favourite id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /favourite/{id} [delete]
func (h *Handler) deleteFavourite(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	_, userId := getUser(c)

	if err := h.services.Favourite.Delete(c, id, userId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
