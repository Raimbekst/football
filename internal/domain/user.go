package domain

type User struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"user_name"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	UserType    string `json:"user_type,omitempty" db:"user_type"`
	Password    string `json:"password,omitempty" db:"password" `
	IsActivated string `json:"is_activated,omitempty" db:"is_activated"`
}

type UserUpdate struct {
	Name        *string `json:"name" db:"user_name"`
	PhoneNumber *string `json:"phone_number" db:"phone_number"`
}

type VerifyUserInput struct {
	Phone     string `json:"phone_number" validate:"required"`
	PhoneCode string `json:"secret_code" validate:"required"`
}

type SetPasswordInput struct {
	CurrentPassword    string `json:"current_password" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required,min=8,max=64" `
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required"`
}

type VerifyPhoneNumberInput struct {
	PhoneNumber string `json:"phone_number"           validate:"required"`
	SecretCode  string `json:"secret_code"            validate:"required"`
}

type ResetPasswordInput struct {
	PhoneNumber        string `json:"phone_number"           validate:"required"`
	NewPassword        string `json:"new_password"           validate:"required,min=8,max=64" `
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type ResetPhoneNumberInput struct {
	PhoneNumber string `json:"phone_number"           validate:"required"`
	SecretCode  string `json:"secret_code"     validate:"required"`
}
