package model

type Event struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Date    string `json:"date"`
}
