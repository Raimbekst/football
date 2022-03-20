package domain

type BuildingImage struct {
	Id         int    `json:"id" db:"id"`
	BuildingId int    `json:"building_id" db:"building_id"`
	Image      string `json:"image" db:"building_image"`
	ManagerId  int    `json:"manager_id,omitempty" `
}
