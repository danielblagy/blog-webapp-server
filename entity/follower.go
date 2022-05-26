package entity

type Follower struct {
	FollowerId int `json:"follower_id" gorm:"not null;index"`
	FollowsId  int `json:"follows_id" gorm:"not null;index"`
}
