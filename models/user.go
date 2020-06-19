package models

//User ...
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Avatar   string `json:"avatar_url"`
}
