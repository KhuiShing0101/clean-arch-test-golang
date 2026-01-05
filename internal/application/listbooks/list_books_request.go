package listbooks

// ListBooksRequest - DTO for list books query
type ListBooksRequest struct {
	Page  int
	Limit int
}

// NewListBooksRequest creates a new request with defaults
func NewListBooksRequest(page, limit int) *ListBooksRequest {
	// Set defaults
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return &ListBooksRequest{
		Page:  page,
		Limit: limit,
	}
}

// GetOffset calculates the SQL offset from page and limit
func (r *ListBooksRequest) GetOffset() int {
	return (r.Page - 1) * r.Limit
}
