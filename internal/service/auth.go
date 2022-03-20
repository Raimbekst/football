package service

import (
	"carWash/internal/domain"
	"carWash/internal/repository"
	"carWash/pkg/auth"
	"carWash/pkg/hash"
	"carWash/pkg/phone"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type UserAuthService struct {
	repos           repository.UserAuth
	hashes          hash.PasswordHashes
	otpPhone        phone.SecretGenerator
	redis           *redis.Client
	ctx             context.Context
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserAuthService(
	repos repository.UserAuth,
	hashes hash.PasswordHashes,
	otpPhone phone.SecretGenerator,
	redis *redis.Client,
	ctx context.Context,
	tokenManager auth.TokenManager,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration) *UserAuthService {
	return &UserAuthService{
		repos:           repos,
		hashes:          hashes,
		otpPhone:        otpPhone,
		redis:           redis,
		ctx:             ctx,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (u *UserAuthService) VerifyExistenceUser(phone string) error {
	_, err := u.repos.VerifyExistenceUser(phone, true)
	if err != nil {
		return fmt.Errorf("service.VerifyExistenceUser: %w", err)
	}
	return nil
}

func (u *UserAuthService) UserSignUp(input SignUpInput) (string, error) {

	if input.Password != input.ConfirmPassword {
		return "", fmt.Errorf("%w", domain.ErrPasswordNotMatch)
	}

	hashedPassword, err := u.hashes.Hash(input.Password)
	if err != nil {
		return "", fmt.Errorf("service.UserSignUp: %w", err)
	}

	user := domain.User{
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		Password:    hashedPassword,
		UserType:    input.UserType,
	}

	code, err := u.SetSecretCode(input.PhoneNumber)
	if err != nil {
		return "", fmt.Errorf("service.UserSignUp: %w", err)
	}

	list, err := u.repos.VerifyExistenceUser(input.PhoneNumber, false)

	fmt.Println(err)
	if err != nil {
		_, err = u.repos.CreateUser(user)
		if err != nil {
			return "", fmt.Errorf("service.UserSignUp: %w", err)
		}

	} else {
		err = u.repos.UpdateUser(user, list.Id)
		if err != nil {
			return "", fmt.Errorf("service.UserSignUp: %w", err)
		}
	}

	return code, nil
}

func (u *UserAuthService) Verify(input domain.VerifyUserInput) error {
	err := u.GetSecretCode(input.PhoneCode, input.Phone)
	if err != nil {
		return fmt.Errorf("service.Verify: %w", err)
	}

	err = u.repos.Verify(input.Phone)
	if err != nil {
		return fmt.Errorf("service.Verify: %w", err)
	}
	return nil
}

func (u *UserAuthService) UserSignIn(user domain.User) (*Tokens, error) {

	hashedPassword, err := u.hashes.Hash(user.Password)

	if err != nil {
		return nil, fmt.Errorf("service.UserSignIn: %w", err)
	}

	input, err := u.repos.SignIn(user.PhoneNumber, hashedPassword)

	if err != nil {
		return nil, fmt.Errorf("service.UserSignIn: %w", err)
	}

	return u.createSession(input.Id, input.UserType)
}

func (u *UserAuthService) SetPassword(id int, input domain.SetPasswordInput) error {

	hashedOldPassword, err := u.hashes.Hash(input.CurrentPassword)

	if err != nil {
		return fmt.Errorf("service.SetPassword: %w", err)
	}

	err = u.repos.VerifyViaPassword(id, hashedOldPassword)

	if err != nil {
		return fmt.Errorf("servcie.SetPassword: %w", err)
	}

	if input.NewPassword != input.ConfirmNewPassword {
		return fmt.Errorf("service.SetPassword: %w", domain.ErrPasswordNotMatch)
	}

	hashedNewPassword, err := u.hashes.Hash(input.NewPassword)

	if err != nil {
		return fmt.Errorf("service.SetPassword: %w", err)
	}

	return u.repos.SetPassword(id, hashedOldPassword, hashedNewPassword)
}

func (u *UserAuthService) ResetPassword(phone string) (string, error) {

	input, err := u.repos.VerifyViaPhoneNumber(phone)
	fmt.Println(input)

	if err != nil {
		return "", fmt.Errorf("service.ResetPassword: %w", err)
	}

	secret, err := u.SetSecretCode(phone)

	if err != nil {
		return "", fmt.Errorf("service.ResetPassword: %w", err)
	}

	if err != nil {
		return "", fmt.Errorf("service.UserSignUp: %w", err)
	}

	return secret, nil

}

func (u *UserAuthService) ResetPasswordConfirm(input domain.ResetPasswordInput) error {
	err := u.GetSecretCode(input.SecretCode, input.PhoneNumber)

	if err != nil {
		return fmt.Errorf("service.ResetPasswordConfirm: %w", err)
	}

	if input.NewPassword != input.ConfirmNewPassword {
		return fmt.Errorf("service.SetPassword: %w", domain.ErrPasswordNotMatch)
	}

	hashedNewPassword, err := u.hashes.Hash(input.NewPassword)

	if err != nil {
		return fmt.Errorf("service.ResetPasswordConfirm: %w", err)
	}

	if err := u.repos.ResetPassword(input.PhoneNumber, hashedNewPassword); err != nil {
		return fmt.Errorf("service.ResetPasswordConfirm : %w", err)
	}
	return nil

}

func (u *UserAuthService) UpdatePhoneNumberVerify(inp domain.User) (string, error) {

	secret, err := u.SetSecretCode(inp.PhoneNumber)

	if err != nil {
		return "", fmt.Errorf("service.ResetPassword: %w", err)
	}
	return secret, nil
}

func (u *UserAuthService) UpdatePhoneNumberConfirm(input domain.ResetPhoneNumberInput, id int) error {
	err := u.GetSecretCode(input.SecretCode, input.PhoneNumber)
	if err != nil {
		return fmt.Errorf("service.UpdatePhoneNumberConfirm: %w", err)
	}
	user := domain.User{PhoneNumber: input.PhoneNumber}

	err = u.repos.UpdateUser(user, id)
	if err != nil {
		return fmt.Errorf("service.UpdatePhoneNumberConfirm: %w", err)
	}
	return nil

}

func (u *UserAuthService) SetSecretCode(phone string) (string, error) {
	secret, err := u.otpPhone.GetRandNum()
	if err != nil {
		return "", fmt.Errorf("service.SetSecretCode: %w", err)
	}
	hashSecret, err := u.hashes.Hash(secret)
	if err != nil {
		return "", fmt.Errorf("service.SetSecretCode: %w", err)
	}
	err = u.redis.Set(u.ctx, phone, hashSecret, 2*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("service.SetSecretCode: %w", err)
	}
	return secret, nil
}

func (u *UserAuthService) GetSecretCode(code, phone string) error {
	hashSecret, err := u.hashes.Hash(code)
	if err != nil {
		return fmt.Errorf("service.GetSecretCode: %w", err)
	}
	val, err := u.redis.Get(u.ctx, phone).Result()
	if err != nil {
		return fmt.Errorf("service.GetSecretCode: %w", err)
	}
	if val != hashSecret {
		return fmt.Errorf("service.GetSecretCode: %w", domain.ErrInvalidSecretCode)
	}
	return nil
}

func (u *UserAuthService) createSession(userId int, userType string) (*Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = u.tokenManager.NewJWT(strconv.Itoa(userId), userType, u.accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("service.createSession.NewJWT: %w", err)

	}

	res.RefreshToken, err = u.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("service.createSession.NewRefreshToken: %w", err)
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	}

	err = u.repos.SetSession(userId, session)

	if err != nil {
		return nil, fmt.Errorf("service.createSession: %w", err)
	}

	return &res, nil
}
