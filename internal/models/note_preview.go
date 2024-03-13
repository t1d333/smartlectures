package models

type NotePreview struct {
	NoteId    int    `json:"noteId"`
	Name      string `json:"name"`
	ParentDir int    `json:"parentDir"`
}
