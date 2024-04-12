package models

type NoteSearchItem struct {
	NoteID        int      `json:"noteId"`
	Name          string   `json:"name"`
	BodyHighlight []string `json:"bodyHighlight"`
	NameHighlight []string `json:"nameHighlight"`
}
