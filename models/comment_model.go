package models

import "time"

type Comment struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	User        User      `json:"owner"`
	Post        string    `json:"post"`
	PublishDate time.Time `json:"publishDate"`
}
