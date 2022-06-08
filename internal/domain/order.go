package domain

type Order struct {
	Id           int            `json:"id"  db:"id"`
	PitchId      int            `json:"pitch_id,omitempty" db:"pitch_id"`
	UserId       int            `json:"user_id,omitempty" db:"user_id"`
	Price        int            `json:"price" db:"price"`
	PitchImage   string         `json:"pitch_image" db:"pitch_image"`
	PitchType    int            `json:"pitch_type" db:"pitch_type"`
	UserName     string         `json:"user_name" db:"first_name"`
	PhoneNumber  string         `json:"phone_number" db:"phone_number"`
	BuildingName string         `json:"building_name" db:"building_name"`
	Address      string         `json:"address" db:"address"`
	OrderDate    float64        `json:"order_date" db:"order_date"`
	Status       int            `json:"status" db:"status"`
	Times        []string       `json:"times,omitempty" `
	ServiceIds   []int          `json:"service_ids,omitempty"`
	CardId       int            `json:"card_id" db:"card_id"`
	ExtraInfo    string         `json:"extra_info"`
	TimeInput    []TimeInput    `json:"collection_times,omitempty"`
	ServiceInput []ServiceInput `json:"collection_service,omitempty"`
}

type ServiceInput struct {
	Id          int    `json:"id,omitempty"  db:"id"`
	Service     string `json:"service_name" db:"service_name"`
	ServiceCost int    `json:"price" db:"price"`
}

type TimeInput struct {
	Id         int    `json:"id,omitempty"  db:"id"`
	WorkTime   string `json:"order_work_time" db:"order_work_time"`
	BuildingId int    `json:"building_id,omitempty" db:"building_id"`
}

type OrderTime struct {
	Id       int    `json:"id,omitempty" db:"id"`
	WorkTime string `json:"work_time" db:"work_time"`
	IsBooked bool   `json:"is_booked" db:"is_booked"`
}

type FilterForOrderTimes struct {
	OrderDate  float64 `json:"order_date" form:"order_date" query:"order_date"`
	BuildingId int     `json:"building_id" form:"building_id" query:"building_id"`
}

type FilterForOrder struct {
	OrderStatus int     `json:"order_status" form:"order_status" query:"order_status" enums:"1,2"`
	OrderDate   float64 `json:"order_date" form:"order_date" query:"order_date"`
}
