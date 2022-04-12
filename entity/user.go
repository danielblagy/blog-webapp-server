package entity

import "encoding/json"

type User struct {
	Id       int       `json:"id" gorm:"primaryKey"`
	Login    string    `json:"login" gorm:"type:varchar(100);uniqueIndex;not null"`
	FullName string    `json:"fullname" gorm:"type:varchar(300);not null"`
	Password string    `json:"password,omitempty" gorm:"type:text;not null"`
	Articles []Article `json:"articles" gorm:"foreignKey:AuthorId"`
}

// remove sensitive imformation from user data in server responses
func (u User) MarshalJSON() ([]byte, error) {
	type user User // prevent recursion
	x := user(u)
	x.Password = "" // we set omitempty in User type, and here we make password empty, so the password propery will be ommited
	return json.Marshal(x)
}
