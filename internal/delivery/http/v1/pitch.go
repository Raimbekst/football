package v1

import (
	"carWash/internal/domain"
	"carWash/pkg/media"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
)

type Pitch struct {
	BuildingId int    `json:"building_id" form:"building_id"`
	Price      int    `json:"price" form:"price"`
	Image      string `json:"image" form:"image"`
	PitchType  int    `json:"pitch_type" form:"pitch_type" enums:"1,2,3"`
	PitchExtra int    `json:"pitch_extra" form:"pitch_extra" enums:"1,2"`
}

type pitchFilter struct {
	BuildingId int `json:"building_id" form:"building_id" query:"building_id"`
}

func (h *Handler) initPitchRoutes(api fiber.Router) {

	partner := api.Group("/pitch")
	{
		partner.Get("/", h.getAllPitch)
		partner.Get("/:id", h.getPitchById)
		admin := partner.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}), isManager)
		{
			admin.Post("", h.createPitch)
			admin.Put("/:id", h.updatePitch)
			admin.Delete("/:id", h.deletePitch)
		}
	}

}

// @Security User_Auth
// @Tags pitch
// @ModuleID createPitch
// @Accept  multipart/form-data
// @Produce  json
// @Param building_id formData int true "building id"
// @Param price  formData int true "price for pitch"
// @Param image formData file  true "pitch image"
// @Param pitch_type formData int true "pitch type" Enums(1,2,3)
// @Param pitch_extra formData int false "pitch extra" Enums(1,2)
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /pitch [post]
func (h *Handler) createPitch(c *fiber.Ctx) error {
	var (
		input Pitch
		err   error
	)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, userId := getUser(c)

	file, _ := c.FormFile("image")
	var img string

	if file != nil {
		img, err = media.GetFileName(c, file)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
	}

	pitch := domain.Pitch{
		Image:      img,
		BuildingId: input.BuildingId,
		Price:      input.Price,
		PitchType:  input.PitchType,
		PitchExtra: input.PitchExtra,
		ManagerId:  userId,
	}

	id, err := h.services.Pitch.Create(c, pitch)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{id})
}

// @Tags pitch
// @Description get all pitch
// @ID get-all-pitch
// @Accept  json
// @Produce  json
// @Param filter query pitchFilter true "pitch filter"
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /pitch [get]
func (h *Handler) getAllPitch(c *fiber.Ctx) error {
	var (
		page   domain.Pagination
		filter pitchFilter
	)

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	list, err := h.services.Pitch.GetAll(c, page, filter.BuildingId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags pitch
// @Description get pitch by id
// @ID get-pitch-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "pitch id"
// @Success 200 {object} domain.Pitch
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /pitch/{id} [get]
func (h *Handler) getPitchById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Pitch.GetById(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags pitch
// @Description  update  pitch
// @ModuleID updatePitch
// @Accept  multipart/form-data
// @Produce  json
// @Param id path string true "pitch id"
// @Param building_id formData int false "building id"
// @Param price  formData int false "price for pitch"
// @Param image formData file  false "pitch image"
// @Param pitch_type formData int false "pitch type" Enums(1,2,3)
// @Param pitch_extra formData int false "pitch extra" Enums(1,2)
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /pitch/{id} [put]
func (h *Handler) updatePitch(c *fiber.Ctx) error {
	var (
		input Pitch
		err   error
	)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	file, _ := c.FormFile("image")

	var img string

	if file != nil {
		img, err = media.GetFileName(c, file)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
	}
	pitch := domain.Pitch{
		Image:      img,
		BuildingId: input.BuildingId,
		Price:      input.Price,
		PitchType:  input.PitchType,
		PitchExtra: input.PitchExtra,
	}

	if err := h.services.Pitch.Update(c, id, pitch); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags pitch
// @Description delete pitch
// @ModuleID deleteCemetery
// @Accept  json
// @Produce  json
// @Param id path string true "pitch id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /pitch/{id} [delete]
func (h *Handler) deletePitch(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Pitch.Delete(c, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
