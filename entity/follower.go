package entity

type Follower struct {
	FollowerId int `json:"follower_id" gorm:"not null;uniqueIndex:idx_follower_follows"`
	FollowsId  int `json:"follows_id" gorm:"not null;uniqueIndex:idx_follower_follows"`
}
