package helper

type Pagination struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

func PaginationResponse(message string, page, limit, totalPages int, totalItems int64, data any) Pagination {
	return Pagination{
		Success:    true,
		Message:    message,
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Data:       data,
	}
}
