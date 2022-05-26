package entity

type Save struct {
	UserId    int `json:"user_id" gorm:"not null;index"`
	ArticleId int `json:"article_id" gorm:"not null;index"`
}
