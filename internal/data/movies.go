package data

import "time"

// All fields are exported (Capital Letter), necessary to
// be visible because of 'enconding/json'

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitzero"`
	Version   int32     `json:"version"`
}
