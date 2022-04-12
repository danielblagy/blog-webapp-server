package entity

type User struct {
	Id       int       `json:"id" gorm:"primaryKey"`
	Login    string    `json:"login" gorm:"type varchar(100);uniqueIndex;not null"`
	FullName string    `json:"fullname" gorm:"not null"`
	Password string    `json:"password" gorm:"not null"`
	Articles []Article `json:"articles" gorm:"foreignKey:AuthorId"`
}
