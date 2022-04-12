package entity

type Article struct {
	Id        int    `json:"id" gorm:"primaryKey"`
	AuthorId  int    `json:"author_id" gorm:"not null"`
	Title     string `json:"title" gorm:"type:varchar(300);not null"`
	Content   string `json:"content" gorm:"type:text;not null"`
	Published bool   `json:"published" gorm:"not null"`
}
