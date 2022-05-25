package domain

type Comment struct {
	Id          int     `json:"id" db:"id"`
	CommentText string  `json:"comment" db:"comment"`
	UserId      int     `json:"user_id,omitempty" db:"user_id"`
	UserName    string  `json:"user_name,omitempty" db:"user_name"`
	BuildingId  int     `json:"building_id" db:"building_id"`
	Grade       float64 `json:"grade,omitempty" db:"grade"`
	PostData    float64 `json:"post_data" db:"post_data"`
}

type Grade struct {
	Id         int     `json:"id" db:"id"`
	UserId     int     `json:"user_id,omitempty" db:"user_id"`
	UserName   string  `json:"user_name,omitempty" db:"user_name"`
	BuildingId int     `json:"building_id" db:"building_id"`
	Grade      float64 `json:"grade,omitempty" db:"grade"`
}
