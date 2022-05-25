package v1

import (
	"carWash/internal/domain"
	"carWash/internal/service"
	"carWash/pkg/validation/validationStructs"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

type signInInput struct {
	PhoneNumber string `json:"phone_number"         validate:"required"`
	Password    string `json:"password"             validate:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type codeResponse struct {
	SecretCode string `json:"secret_code,omitempty"`
}

type UpdateNumberInput struct {
	NewPhoneNumber string `json:"new_phone_number" validate:"required"`
}

type PhoneNumberInput struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

func (h *Handler) initUserRoutes(api fiber.Router) {
	auth := api.Group("/auth")
	{
		auth.Post("/user/sign-up", h.userSignUp)
		auth.Post("/manager/sign-up", h.managerSignUp)
		auth.Post("verify", h.verifyUser)
		auth.Post("sign-in", h.userSignIn)
		auth.Post("reset-password", h.resetPasswordVerify)
		auth.Post("reset-password-verify-phone-number", h.verifyPhoneNumberForResetPassword)
		auth.Post("reset-password-confirm", h.resetPasswordConfirm)

		users := auth.Group("").Use(jwtware.New(
			jwtware.Config{
				SigningKey: []byte(h.signingKey),
			}), isUser)
		{
			users.Get("user", h.getUser)
			users.Put("user", h.updateUser)
			users.Post("set-password", h.userSetPassword)
			users.Post("set-phone-number-verify", h.updateNumber)
			users.Post("set-phone-number-confirm", h.updateNumberConfirm)
		}
	}
}

type UserSignUpInput struct {
	Name            string `json:"name"`
	PhoneNumber     string `json:"phone_number"       validate:"required,e164" `
	Password        string `json:"password,omitempty" validate:"required,min=8,max=64" `
	ConfirmPassword string `json:"confirm_password"   validate:"required"`
}

// @Tags auth
// @Description create member account
// @ModuleID managerSignUp
// @Accept json
// @Produce  json
// @Param data body UserSignUpInput true "manager sign-up"
// @Success 201 {object} codeResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/manager/sign-up [post]
func (h *Handler) managerSignUp(c *fiber.Ctx) error {
	var input UserSignUpInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}
	err := h.services.UserAuth.VerifyExistenceUser(input.PhoneNumber)

	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: domain.ErrUserAlreadyExist.Error()})
	}

	secret, err := h.services.UserAuth.UserSignUp(service.SignUpInput{
		Name:            input.Name,
		PhoneNumber:     input.PhoneNumber,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
		UserType:        manager,
	})

	if err != nil {
		if errors.Is(err, domain.ErrPasswordNotMatch) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(codeResponse{SecretCode: secret})
}

// @Tags auth
// @Description create member account
// @ModuleID userSignUp
// @Accept json
// @Produce  json
// @Param data body UserSignUpInput true "user sign-up"
// @Success 201 {object} codeResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/user/sign-up [post]
func (h *Handler) userSignUp(c *fiber.Ctx) error {
	var input UserSignUpInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}
	err := h.services.UserAuth.VerifyExistenceUser(input.PhoneNumber)

	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: domain.ErrUserAlreadyExist.Error()})
	}

	secret, err := h.services.UserAuth.UserSignUp(service.SignUpInput{
		Name:            input.Name,
		PhoneNumber:     input.PhoneNumber,
		Password:        input.Password,
		ConfirmPassword: input.ConfirmPassword,
		UserType:        user,
	})

	if err != nil {
		if errors.Is(err, domain.ErrPasswordNotMatch) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(codeResponse{SecretCode: secret})
}

// @Tags auth
// @Description verify user phone number
// @ModuleID verifyUser
// @Accept json
// @Produce json
// @Param data body domain.VerifyUserInput true "user verify"
// @Success 201 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/verify [post]
func (h *Handler) verifyUser(c *fiber.Ctx) error {
	var input domain.VerifyUserInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	err := h.services.UserAuth.Verify(input)

	if err != nil {
		if errors.Is(err, domain.ErrInvalidSecretCode) || errors.Is(err, domain.ErrUserAlreadyExist) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(okResponse{Message: "OK"})
}

// @Tags auth
// @Description user sign in
// @ModuleID userSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/sign-in [post]
func (h *Handler) userSignIn(c *fiber.Ctx) error {
	var input signInInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}
	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	user := domain.User{
		PhoneNumber: input.PhoneNumber,
		Password:    input.Password,
	}

	res, err := h.services.UserAuth.UserSignIn(user)

	if err != nil {

		if errors.Is(err, domain.ErrUserDoesNotExist) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tokenResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

// @Tags auth
// @Security User_Auth
// @Description users change password
// @ModuleID userSetPassword
// @Accept  json
// @Produce  json
// @Param input body domain.SetPasswordInput true "change password"
// @Success 201 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/set-password [post]
func (h *Handler) userSetPassword(c *fiber.Ctx) error {

	var input domain.SetPasswordInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	_, id := getUser(c)

	err := h.services.UserAuth.SetPassword(id, input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPassword) || errors.Is(err, domain.ErrPasswordNotMatch) {
			return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "Ok"})
}

// @Tags auth
// @Description verify email for reset password
// @ModuleID resetPasswordVerify
// @Accept json
// @Produce json
// @Param data body PhoneNumberInput true "enter email"
// @Success 201 {object} codeResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/reset-password [post]
func (h *Handler) resetPasswordVerify(c *fiber.Ctx) error {
	var input PhoneNumberInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	secret, err := h.services.UserAuth.ResetPassword(input.PhoneNumber)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotRegistered) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(codeResponse{secret})
}

// @Tags auth
// @Description verify secret code for reset password
// @ModuleID verifyPhoneNumberForResetPassword
// @Accept json
// @Produce json
// @Param data body domain.VerifyPhoneNumberInput true "verify phone number"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/reset-password-verify-phone-number [post]
func (h *Handler) verifyPhoneNumberForResetPassword(c *fiber.Ctx) error {
	var input domain.VerifyPhoneNumberInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	err := h.services.UserAuth.VerifyPhoneNumber(input)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, domain.ErrInvalidSecretCode) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Tags auth
// @Description reset password
// @ModuleID resetPasswordConfirm
// @Accept json
// @Produce json
// @Param data body domain.ResetPasswordInput true "reset password"
// @Success 200 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/reset-password-confirm [post]
func (h *Handler) resetPasswordConfirm(c *fiber.Ctx) error {
	var input domain.ResetPasswordInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	ok, errs := validationStructs.ValidateStruct(input)

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	err := h.services.UserAuth.ResetPasswordConfirm(input)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, domain.ErrInvalidSecretCode) || errors.Is(err, domain.ErrPasswordNotMatch) {
			return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})

}

// @Security User_Auth
// @Tags auth
// @Description update phone number verify
// @ModuleID updateNumber
// @Accept json
// @Produce json
// @Param data body UpdateNumberInput true "phone number verify input"
// @Success 201 {object} codeResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/set-phone-number-verify [post]
func (h *Handler) updateNumber(c *fiber.Ctx) error {
	var input UpdateNumberInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	inp := domain.User{
		PhoneNumber: input.NewPhoneNumber,
	}

	secret, err := h.services.UserAuth.UpdatePhoneNumberVerify(inp)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(codeResponse{SecretCode: secret})
}

// @Security User_Auth
// @Tags auth
// @Description update phone number execute
// @ModuleID updateNumberConfirm
// @Accept json
// @Produce json
// @Param data body domain.ResetPhoneNumberInput true "phone number confirm"
// @Success 201 {object} codeResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/set-phone-number-confirm [post]
func (h *Handler) updateNumberConfirm(c *fiber.Ctx) error {
	var input domain.ResetPhoneNumberInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, id := getUser(c)

	err := h.services.UserAuth.UpdatePhoneNumberConfirm(input, id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(okResponse{Message: "OK"})
}

// @Tags auth
// @Security User_Auth
// @Description get user info
// @ModuleID getUser
// @Accept  json
// @Produce  json
// @Success 201 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/user [get]
func (h *Handler) getUser(c *fiber.Ctx) error {

	_, id := getUser(c)

	list, err := h.services.UserAuth.GetUserInfo(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(list)
}

// @Tags auth
// @Security User_Auth
// @Description update user info
// @ModuleID getUser
// @Accept  json
// @Produce  json
// @Param data body domain.UserUpdate true "a user update info"
// @Success 201 {object} okResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /auth/user [put]
func (h *Handler) updateUser(c *fiber.Ctx) error {

	var inp domain.UserUpdate

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	_, id := getUser(c)

	err := h.services.UserAuth.UpdateUserInfo(inp, id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON("OK")
}
