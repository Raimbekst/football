package domain

import "time"

type Building struct {
	Id          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"building_name"`
	Address     string    `json:"address" db:"address"`
	Instagram   string    `json:"instagram" db:"instagram"`
	Description string    `json:"description" db:"description"`
	ManagerId   int       `json:"manager_id" db:"manager_id"`
	WorkTime    int       `json:"work_time" db:"work_time"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
}

type UserInfo struct {
	Id   int
	Type string
}
