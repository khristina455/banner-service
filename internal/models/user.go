package models

type User struct {
	UserID   int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	TagID    int    `json:"tag_id"`
}
