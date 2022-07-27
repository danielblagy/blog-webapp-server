package entity

type Save struct {
	UserId    int `json:"user_id" gorm:"not null;uniqueIndex:idx_user_article"`
	ArticleId int `json:"article_id" gorm:"not null;uniqueIndex:idx_user_article"`
}
