package models

import "time"

type Note struct {
	Name        string    `json:"name"`
	NoteId      int       `json:"noteId,omitempty"`
	ParentDir   int       `json:"parentDir"`
	Body        string    `json:"body"`
	UserId      int       `json:"userId"`
	RepeatedNum int       `json:"-"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	LastUpdate  time.Time `json:"lastUpdate,omitempty"`
}
