package domain

type Card struct {
	Id         int    `json:"id" db:"id"`
	Cvv        int    `json:"cvv" db:"cvv"`
	UserId     int    `json:"user_id" db:"user_id"`
	FullName   string `json:"full_name" db:"full_name"`
	FullNumber string `json:"full_number" db:"full_number"`
}
