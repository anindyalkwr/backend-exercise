package models

import "time"

type Post struct {
	ID          string    `json:"id"`
	Image       string    `json:"image"`
	Likes       int       `json:"likes"`
	Tags        []string  `json:"tags"`
	Text        string    `json:"text"`
	PublishDate time.Time `json:"publishDate"`
	User        User      `json:"owner"`
}
