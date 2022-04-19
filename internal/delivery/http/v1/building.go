package v1

import (
	"carWash/internal/domain"
	"carWash/pkg/validation/validationStructs"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"strconv"
)

type Building struct {
	Name        string `json:"name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Instagram   string `json:"instagram"`
	Description string `json:"description" validate:"required"`
	WorkTime    int    `json:"work_time"   validate:"required" enums:"1,2" examples:"1" `
	StartTime   string `json:"start_time"  validate:"required"`
	EndTime     string `json:"end_time"    validate:"required"`
	Longtitude  string `json:"longtitude"  validate:"required"`
	Latitude    string `json:"latitude"    validate:"required"`
}

type UpdateBuilding struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Instagram   string `json:"instagram"`
	Description string `json:"description"`
	WorkTime    int    `json:"work_time" enums:"1,2" examples:"1"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Longtitude  string `json:"longtitude"`
	Latitude    string `json:"latitude"`
}

func (h *Handler) initBuildingRoutes(api fiber.Router) {
	partner := api.Group("/building")
	{
		h.initBuildingImageRoutes(partner)
		partner.Get("/", h.getAllBuildings)
		partner.Get("/:id", h.getBuildingById)
		admin := partner.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}), isManager)
		{
			admin.Post("", h.createBuilding)
			admin.Put("/:id", h.updateBuilding)
			admin.Delete("/:id", h.deleteBuilding)
		}
	}
}

// @Security User_Auth
// @Tags building
// @ModuleID createBuilding
// @Accept json
// @Produce  json
// @Param data body Building true "building create input"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building [post]
func (h *Handler) createBuilding(c *fiber.Ctx) error {
	var (
		input Building
		err   error
	)

	_, userId := getUser(c)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	building := domain.Building{
		Name:        input.Name,
		Address:     input.Address,
		Instagram:   input.Instagram,
		ManagerId:   userId,
		Description: input.Description,
		WorkTime:    input.WorkTime,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		Latitude:    input.Latitude,
		Longtitude:  input.Longtitude,
	}

	id, err := h.services.Building.Create(c, building)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{id})
}

// @Tags building
// @Description get all buildings
// @ID get-all-buildings
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Param filter query domain.FilterForBuilding true "filter for building"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building [get]
func (h *Handler) getAllBuildings(c *fiber.Ctx) error {
	var (
		page   domain.Pagination
		filter domain.FilterForBuilding
		err    error
	)

	header := string(c.Request().Header.Peek("Authorization"))

	userId, userType, err := h.userIdentity(header)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	info := domain.UserInfo{
		Id:   userId,
		Type: userType,
	}

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	fmt.Println(info)

	list, err := h.services.Building.GetAll(c, page, info, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags building
// @Description get building by id
// @ID get-building-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "building id"
// @Success 200 {object} domain.Building
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building/{id} [get]
func (h *Handler) getBuildingById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Building.GetById(c, id)

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(list)
}

// @Security User_Auth
// @Tags building
// @Description  update  building
// @ModuleID updateBuilding
// @Accept json
// @Produce  json
// @Param id path string true "building id"
// @Param data body UpdateBuilding true "building input"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building/{id} [put]
func (h *Handler) updateBuilding(c *fiber.Ctx) error {
	var (
		input UpdateBuilding
		err   error
	)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, userId := getUser(c)

	building := domain.Building{
		Name:        input.Name,
		Address:     input.Address,
		Instagram:   input.Instagram,
		ManagerId:   userId,
		Description: input.Description,
		WorkTime:    input.WorkTime,
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
	}

	if err := h.services.Building.Update(c, id, building); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags building
// @Description delete building
// @ModuleID deleteBuilding
// @Accept  json
// @Produce  json
// @Param id path string true "building id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building/{id} [delete]
func (h *Handler) deleteBuilding(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.Building.Delete(c, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
