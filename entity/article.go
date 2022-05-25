package entity

import "time"

type Article struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	AuthorId  int       `json:"author_id" gorm:"not null"`
	Title     string    `json:"title" gorm:"type:varchar(300);not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Published bool      `json:"published" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Author    User      `json:"author" gorm:"-"`
}

type EditableArticleData struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}
