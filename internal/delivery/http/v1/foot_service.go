package v1

import (
	"carWash/internal/domain"
	"carWash/pkg/validation/validationStructs"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type FootService struct {
	ServiceName string `json:"service_name" validate:"required"`
	Price       int    `json:"price" validate:"required"`
}
type UpdateFootService struct {
	ServiceName string `json:"service_name" `
	Price       int    `json:"price"`
}

func (h *Handler) initFootServiceRoutes(api fiber.Router) {

	partner := api.Group("/service")
	{
		partner.Get("/", h.getAllFootService)
		partner.Get("/:id", h.getFootServiceById)
		partner.Post("", h.createFootService)
		partner.Put("/:id", h.updateFootService)
		partner.Delete("/:id", h.deleteFootService)
	}
}

// @Tags service
// @ModuleID createFootService
// @Accept  json
// @Produce  json
// @Param data body FootService true "service input"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /service [post]
func (h *Handler) createFootService(c *fiber.Ctx) error {
	var (
		input FootService
		err   error
	)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, mess := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(mess)
	}

	service := domain.FootService{
		Price:       input.Price,
		ServiceName: input.ServiceName,
	}

	id, err := h.services.FootService.Create(c, service)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{id})
}

// @Tags service
// @Description get all service
// @ID get-all-service
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /service [get]
func (h *Handler) getAllFootService(c *fiber.Ctx) error {
	var (
		page domain.Pagination
	)

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	list, err := h.services.FootService.GetAll(c, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags service
// @Description get service by id
// @ID get-service-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "service id"
// @Success 200 {object} domain.FootService
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /service/{id} [get]
func (h *Handler) getFootServiceById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.FootService.GetById(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags service
// @Description  update  service
// @ModuleID updateFootService
// @Accept  json
// @Produce  json
// @Param id path string true "service id"
// @Param data body UpdateFootService true "foot service input"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /service/{id} [put]
func (h *Handler) updateFootService(c *fiber.Ctx) error {
	var (
		input UpdateFootService
		err   error
	)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	id, err := strconv.Atoi(c.Params("id"))

	service := domain.FootService{
		Price:       input.Price,
		ServiceName: input.ServiceName,
	}

	if err := h.services.FootService.Update(c, id, service); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Tags service
// @Description delete service
// @ModuleID deleteFootService
// @Accept  json
// @Produce  json
// @Param id path string true "service id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /service/{id} [delete]
func (h *Handler) deleteFootService(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.FootService.Delete(c, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
