package resources

// School is the frontend representation of a single school.
type School struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required""`
}

// Schools is the frontend representation of a paginated collection of schools.
type Schools struct {
	Schools []School `json:"schools"`
	Meta    Meta     `json:"meta"`
	Links   Links    `json:"links"`
}

// Meta is the frontend reprensetation of the metadata associated with a
// collection including the total number of schools in the database.
type Meta struct {
	Total int `json:"total"`
}

// Links is the frontend representation of a set of URLs that represent
// different pages of interest within a paginated collection.
type Links struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}
