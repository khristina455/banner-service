package models

type User struct {
	UserId   int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
	TagId    int    `json:"tag_id"`
}
