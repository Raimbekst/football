package domain

type Order struct {
	Id           int     `json:"id"  db:"id"`
	PitchId      int     `json:"pitch_id,omitempty" db:"pitch_id"`
	UserId       int     `json:"user_id,omitempty" db:"user_id"`
	Price        int     `json:"price" db:"price"`
	PitchImage   string  `json:"pitch_image" db:"pitch_image"`
	PitchType    int     `json:"pitch_type" db:"pitch_type"`
	UserName     string  `json:"user_name" db:"user_name"`
	PhoneNumber  string  `json:"phone_number" db:"phone_number"`
	BuildingName string  `json:"building_name" db:"building_name"`
	Address      string  `json:"address" db:"address"`
	OrderDate    float64 `json:"order_date" db:"order_date"`
	StartTime    string  `json:"start_time" db:"start_time"`
	EndTime      string  `json:"end_time" db:"end_time"`
	Status       int     `json:"status" db:"status"`
}
