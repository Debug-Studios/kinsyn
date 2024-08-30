package commons

import "time"

type Highlight struct {
	BookTitle         string    `json:"book_title"`
	BookAuthor        string    `json:"book_author"`
	BookLocationStart int       `json:"book_location_start"`
	BookLocationEnd   int       `json:"book_location_end"`
	CreatedAt         time.Time `json:"created_at"`
	Content           string    `json:"content"`
}
