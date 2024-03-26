package models

type Dir struct {
	Name        string `json:"name"`
	DirId       int    `json:"dirId"`
	UserId      int    `json:"userId"`
	RepeatedNum int    `json:"-"`
	ParentDir   int    `json:"parentDir"`
	Subdirs     []Dir  `json:"subdirs"`
}
