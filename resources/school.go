package resources

type School struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required""`
}

type Schools struct {
	Schools []School `json:"schools"`
	// TotalCount int      `json:"total_count"`
	// Next       string   `json:"next"`
	// Prev       string   `json:"prev"`
}
