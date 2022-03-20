package domain

type Pitch struct {
	Id         int    `json:"id" db:"id"`
	BuildingId int    `json:"building_id" db:"building_id"`
	Price      int    `json:"price" db:"price"`
	Image      string `json:"image" db:"pitch_image"`
	ManagerId  int    `json:"manager_id,omitempty" `
}
