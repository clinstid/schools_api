package resources

type School struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required""`
}

type Schools struct {
	Schools []School `json:"schools"`
	Meta    Meta     `json:"meta"`
	Links   Links    `json:"links"`
}

type Meta struct {
	Total int `json:"total"`
}

type Links struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}
