package v1

import (
	"carWash/internal/domain"
	"carWash/pkg/validation/validationStructs"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func (h *Handler) initOrderRoutes(api fiber.Router) {
	order := api.Group("/order")
	{
		order.Get("/", h.getAllOrders)
		order.Get("/times", h.getOrderForCreateOrder)

		admin := order.Group("/", jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey)}), isUser)
		{
			admin.Post("/", h.createOrder)
		}
	}
}

type order struct {
	PitchId     int      `json:"pitch_id"   validate:"required"`
	OrderDate   float64  `json:"order_date" validate:"required"`
	Times       []string `json:"times" validate:"required"`
	ServiceIds  []int    `json:"service_ids"`
	CardId      int      `json:"card_id"`
	FirstName   string   `json:"first_name"`
	PhoneNumber string   `json:"phone_number" validate:"required"`
	ExtraInfo   string   `json:"extra_info"`
}

// @Security User_Auth
// @Tags orders
// @ModuleID createOrder
// @Accept json
// @Produce  json
// @Param data body order true "order create input"
// @Success 201 {object} idResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /order [post]
func (h *Handler) createOrder(c *fiber.Ctx) error {

	var (
		inp order
	)

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, mess := validationStructs.ValidateStruct(inp)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(mess)
	}

	_, userId := getUser(c)

	input := domain.Order{
		PitchId:     inp.PitchId,
		UserId:      userId,
		OrderDate:   inp.OrderDate,
		Times:       inp.Times,
		ServiceIds:  inp.ServiceIds,
		CardId:      inp.CardId,
		UserName:    inp.FirstName,
		PhoneNumber: inp.PhoneNumber,
		ExtraInfo:   inp.ExtraInfo,
		Status:      reserved,
	}

	id, err := h.services.Order.Create(c, input)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(idResponse{ID: id})

}

// @Tags orders
// @Description gets all orders
// @ID get-all-orders
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Param filter query domain.FilterForOrder true "filter for orders"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /order [get]
func (h *Handler) getAllOrders(c *fiber.Ctx) error {
	var (
		page   domain.Pagination
		filter domain.FilterForOrder
	)

	if err := c.QueryParser(&page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	header := string(c.Request().Header.Peek("Authorization"))

	userId, userType, err := h.userIdentity(header)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	info := domain.UserInfo{
		Id:   userId,
		Type: userType,
	}

	list, err := h.services.Order.GetAll(c, page, info, filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags orders
// @Description gets all orders time
// @ID get-all-order-times
// @Accept  json
// @Produce  json
// @Param array query domain.Pagination  true "A page info"
// @Param filter query domain.FilterForOrderTimes true "filter for order times"
// @Success 200 {object} domain.GetAllResponses
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /order/times [get]
func (h *Handler) getOrderForCreateOrder(c *fiber.Ctx) error {
	var (
		filter domain.FilterForOrderTimes
	)

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	list, err := h.services.Order.GetAllBookTime(c, filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response{err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}
