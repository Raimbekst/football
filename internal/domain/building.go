package domain

type Building struct {
	Id              int     `json:"id,omitempty" db:"id"`
	Name            string  `json:"name" db:"building_name"`
	Address         string  `json:"address" db:"address"`
	PhoneNumber     string  `json:"phone_number" db:"phone_number"`
	Instagram       string  `json:"instagram" db:"instagram"`
	Description     string  `json:"description" db:"description"`
	BuildingImage   string  `json:"building_image" db:"building_image"`
	ManagerId       int     `json:"manager_id,omitempty" db:"manager_id"`
	WorkTime        int     `json:"work_time,omitempty" db:"work_time"`
	StartTime       int     `json:"start-time,omitempty"`
	StartTimeString string  `json:"start_time,omitempty" db:"start_time"`
	EndTimeString   string  `json:"end_time,omitempty" db:"end_time"`
	EndTime         int     `json:"end-time,omitempty"`
	MinPrice        int     `json:"price" db:"min_price"`
	Longtitude      string  `json:"longitude,omitempty"  db:"longtitude"`
	Latitude        string  `json:"latitude,omitempty"    db:"latitude"`
	IsFavourite     bool    `json:"is_favourite" db:"is_favourite"`
	Grade           float64 `json:"grade" db:"grade"`
	Favourite       `json:"favourite,omitempty" db:"f"`
}

type UserInfo struct {
	Id   int
	Type string
}

type FilterForBuilding struct {
	PitchType  int  `json:"pitch_type" form:"pitch_type" query:"pitch_type"  enums:"1,2,3"`
	PitchExtra int  `json:"pitch_extra" form:"pitch_extra" query:"pitch_extra" enums:"1,2"`
	StartCost  *int `json:"start_cost" form:"start_cost" query:"start_cost"`
	EndCost    *int `json:"end_cost" form:"end_cost" query:"end_cost"`
}
