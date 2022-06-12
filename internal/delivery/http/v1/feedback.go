package v1

import (
	"carWash/internal/domain"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type Feedback struct {
	Text string `json:"text" `
}

type Notification struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *Handler) initNotificationRoutes(api fiber.Router) {

	partner := api.Group("/notification")
	{
		partner.Post("/", h.createNotification)
		partner.Get("/", h.getAllNotifications)
	}
}

func (h *Handler) initFeedbackRoutes(api fiber.Router) {

	partner := api.Group("/feedback")
	{
		partner.Get("/", h.getAllFeedbacks)

		admin := partner.Group("", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}), isUser)
		{
			admin.Post("", h.createFeedback)
		}
	}
}

// @Security User_Auth
// @Tags feedback
// @ModuleID createFeedback
// @Accept  json
// @Produce  json
// @Param data body Feedback true "feedback input"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /feedback [post]
func (h *Handler) createFeedback(c *fiber.Ctx) error {
	var (
		input Feedback
		err   error
	)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, userId := getUser(c)

	inp := domain.Feedback{
		Text: input.Text,
	}
	id, err := h.services.Feedback.Create(c, inp, userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{id})
}

// @Tags feedback
// @Description gets all feedbacks
// @ID get-all-feedback
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /feedback [get]
func (h *Handler) getAllFeedbacks(c *fiber.Ctx) error {
	var (
		page domain.Pagination
	)

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	list, err := h.services.Feedback.GetAll(c, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

/////////////////////////////////////////////////////////////////////

// @Security User_Auth
// @Tags notification
// @ModuleID createNotification
// @Accept  json
// @Produce  json
// @Param data body Notification true "notification input"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification [post]
func (h *Handler) createNotification(c *fiber.Ctx) error {
	var (
		input Notification
		err   error
	)

	if err = c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	inp := domain.Notification{
		Title:   input.Title,
		Content: input.Content,
	}
	id, err := h.services.Feedback.CreateNoty(c, inp)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{id})
}

// @Tags notification
// @Description gets all notifications
// @ID get-all-notification
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /notification [get]
func (h *Handler) getAllNotifications(c *fiber.Ctx) error {
	var (
		page domain.Pagination
	)

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	list, err := h.services.Feedback.GetAllNoty(c, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}
