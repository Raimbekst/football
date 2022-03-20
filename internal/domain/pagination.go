package domain

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type PaginationPage struct {
	Page  int `json:"page"`
	Pages int `json:"pages"`
	Count int `json:"count"`
}

type GetAllResponses struct {
	Data     interface{}    `json:"data"`
	PageInfo PaginationPage `json:"page_info"`
}
