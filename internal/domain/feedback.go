package domain

type Feedback struct {
	Id          int    `json:"id" db:"id"`
	Text        string `json:"text" db:"text"`
	UserId      int    `json:"user_id" db:"user_id"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	UserName    string `json:"user_name" db:"user_name"`
}

type Notification struct {
	Id      int    `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}
