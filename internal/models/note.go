package models

type Note struct {
	Name       string `json:"name"`
	NoteId     int    `json:"noteId"`
	ParentDir  int    `json:"parentDir"`
	Body       string `json:"body"`
	UserId     int    `json:"userId"`
	CreatedAt  string `json:"createdAt"`
	LastUpdate string `json:"lastUpdate"`
}
