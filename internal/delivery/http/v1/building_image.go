package v1

import (
	"carWash/internal/domain"
	"carWash/pkg/media"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"mime/multipart"
	"strconv"
)

type BuildingImage struct {
	BuildingId int                   `json:"building_id" form:"building_id"`
	Image      *multipart.FileHeader `json:"image"       form:"image"`
}

type buildingImageFilter struct {
	BuildingId int `json:"building_id" db:"building_id"`
}

func (h *Handler) initBuildingImageRoutes(api fiber.Router) {
	partner := api.Group("/image")
	{
		partner.Get("/", h.getAllBuildingImages)
		partner.Get("/:id", h.getBuildingImageById)
		admin := partner.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}))
		{
			admin.Post("", h.createBuildingImage)
			admin.Put("/:id", h.updateBuildingImage)
			admin.Delete("/:id", h.deleteBuildingImage)
		}
	}
}

// @Security User_Auth
// @Tags building
// @ModuleID createBuildingImage
// @Accept  multipart/form-data
// @Produce  json
// @Param building_id formData int true "building id"
// @Param image  formData file true "image of building"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building/image [post]
func (h *Handler) createBuildingImage(c *fiber.Ctx) error {
	var (
		input BuildingImage
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

	inp := domain.BuildingImage{
		BuildingId: input.BuildingId,
		Image:      img,
		ManagerId:  userId,
	}

	id, err := h.services.BuildingImage.Create(c, inp)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{id})
}

// @Tags building
// @Description get all buildings
// @ID get-all-building
// @Accept  json
// @Produce  json
// @Param filter query buildingImageFilter true "building filter input"
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building/image [get]
func (h *Handler) getAllBuildingImages(c *fiber.Ctx) error {
	var (
		filter buildingImageFilter
		page   domain.Pagination
	)

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	list, err := h.services.BuildingImage.GetAll(c, page, filter.BuildingId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags building
// @Description get building by id
// @ID get-building-image-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "building id"
// @Success 200 {object} domain.BuildingImage
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building/image/{id} [get]
func (h *Handler) getBuildingImageById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.BuildingImage.GetById(c, id)

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
// @ModuleID updateBuildingImage
// @Accept  multipart/form-data
// @Produce  json
// @Param id path string true "building id"
// @Param building_id formData int false "building id"
// @Param image  formData file false "image of building"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building/image/{id} [put]
func (h *Handler) updateBuildingImage(c *fiber.Ctx) error {
	var (
		input BuildingImage
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

	inp := domain.BuildingImage{
		BuildingId: input.BuildingId,
		Image:      img,
	}

	if err := h.services.BuildingImage.Update(c, id, inp); err != nil {
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
// @ModuleID deleteBuildingImage
// @Accept  json
// @Produce  json
// @Param id path string true "building id"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /building/image/{id} [delete]
func (h *Handler) deleteBuildingImage(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	if err := h.services.BuildingImage.Delete(c, id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}
