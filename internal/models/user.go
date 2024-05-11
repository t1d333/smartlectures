package models

type User struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Surname  string `json:"surname"`
	Name     string `json:"name"`
	Password []byte `json:"-"`
}
