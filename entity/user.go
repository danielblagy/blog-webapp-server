package entity

type User struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
}
