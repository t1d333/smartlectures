package models

type Dir struct {
	Name      string `json:"name"`
	DirId     int    `json:"dirId"`
	UserId    int    `json:"userId"`
	ParentDir int    `json:"parentDir"`
}
