package entity

type Article struct {
	Id        int `json:"id" gorm:"primaryKey"`
	AuthorId  int
	Title     string
	Content   string
	Published bool
}
