package v1

import (
	"carWash/internal/domain"
	"carWash/pkg/validation/validationStructs"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
)

type Card struct {
	Cvv        int    `json:"cvv"         validate:"required"`
	FullName   string `json:"full_name"   validate:"required"`
	FullNumber string `json:"full_number" validate:"required"`
}

type UpdateCard struct {
	Cvv        int    `json:"cvv"         validate:"required"`
	FullName   string `json:"full_name"   validate:"required"`
	FullNumber string `json:"full_number" validate:"required"`
}

func (h *Handler) initCardRoutes(api fiber.Router) {

	card := api.Group("/card")
	{
		admin := card.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}), isUser)
		{
			admin.Get("/", h.getAllCard)
			admin.Get("/:id", h.getCardById)
			admin.Post("", h.createCard)
			admin.Put("/:id", h.updateCard)
			admin.Delete("/:id", h.deleteCard)
		}
	}

}

// @Security User_Auth
// @Tags card
// @ModuleID createCard
// @Accept  json
// @Produce  json
// @Param data body Card true "card input"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /card [post]
func (h *Handler) createCard(c *fiber.Ctx) error {
	var (
		input Card
		err   error
	)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, mess := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(mess)
	}

	_, userId := getUser(c)

	card := domain.Card{
		Cvv:        input.Cvv,
		UserId:     userId,
		FullName:   input.FullName,
		FullNumber: input.FullNumber,
	}

	id, err := h.services.Card.Create(c, card)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{id})
}

// @Security User_Auth
// @Tags card
// @Description get all card
// @ID get-all-card
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /card [get]
func (h *Handler) getAllCard(c *fiber.Ctx) error {
	var (
		page domain.Pagination
	)

	_, userId := getUser(c)

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	list, err := h.services.Card.GetAll(c, page, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags card
// @Description get card by id
// @ID get-card-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "card id"
// @Success 200 {object} domain.Card
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /card/{id} [get]
func (h *Handler) getCardById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	_, userId := getUser(c)

	list, err := h.services.Card.GetById(c, id, userId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags card
// @Description  update  card
// @ModuleID updateCard
// @Accept  json
// @Produce  json
// @Param id path string true "card id"
// @Param data body UpdateCard true "foot card input"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /card/{id} [put]
func (h *Handler) updateCard(c *fiber.Ctx) error {
	var (
		input UpdateCard
		err   error
	)

	_, userId := getUser(c)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	id, err := strconv.Atoi(c.Params("id"))

	card := domain.Card{
		Cvv:        input.Cvv,
		UserId:     userId,
		FullName:   input.FullName,
		FullNumber: input.FullNumber,
	}

	if err := h.services.Card.Update(c, id, card); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags card
// @Description delete card
// @ModuleID deleteCard
// @Accept  json
// @Produce  json
// @Param id path string true "card id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /card/{id} [delete]
func (h *Handler) deleteCard(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, userId := getUser(c)

	if err := h.services.Card.Delete(c, id, userId); err != nil {

		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
