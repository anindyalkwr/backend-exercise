package models

type User struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Picture   string `json:"picture"`
}
