package models

type Snippet struct {
	SnippetID   int    `json:"snippetId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
	UserId      int    `json:"userId"`
}
