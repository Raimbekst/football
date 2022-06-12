package v1

import (
	"carWash/internal/domain"
	"carWash/pkg/validation/validationStructs"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Comment struct {
	CommentText string  `json:"comment" validate:"required"`
	BuildingId  int     `json:"building_id" validate:"required"`
	Grade       float64 `json:"grade" validate:"required"`
}

func (h *Handler) initCommentRoutes(api fiber.Router) {

	partner := api.Group("/comment")
	{
		partner.Get("/", h.getAllComments)
		admin := partner.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}), isUser)
		{
			admin.Post("", h.createComment)
		}
	}
}

// @Security User_Auth
// @Tags comment
// @ModuleID createComment
// @Accept json
// @Produce json
// @Param data body Comment true "comment info"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /comment [post]
func (h *Handler) createComment(c *fiber.Ctx) error {
	var input Comment

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, userId := getUser(c)

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	inp := domain.Comment{
		CommentText: input.CommentText,
		UserId:      userId,
		BuildingId:  input.BuildingId,
		Grade:       input.Grade,
	}

	id, err := h.services.Comment.Create(c, inp)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})
}

// @Tags comment
// @Description get all comments
// @ID get-all-comments
// @Accept json
// @Produce json
// @Param array query domain.Pagination true "a page info"
// @Param filter query pitchFilter true "comment filter"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /comment [get]
func (h *Handler) getAllComments(c *fiber.Ctx) error {
	var page domain.Pagination
	var filter pitchFilter

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	list, err := h.services.Comment.GetAll(c, page, filter.BuildingId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)

}
