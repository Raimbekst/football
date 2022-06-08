package domain

type FootService struct {
	Id          int    `json:"id" db:"id"`
	ServiceName string `json:"service_name" db:"service_name"`
	Price       int    `json:"price" db:"price"`
}
