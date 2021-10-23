package models

type Song struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`

	// Date     time.Time
	// Duration time.Time
}

type Artist struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}
