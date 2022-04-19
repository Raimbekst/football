package domain

type Favourite struct {
	Id       int `json:"id" db:"id"`
	Building `json:"building" db:"b"`
}
