package model

import "fmt"

type FollowStruct struct {
	BaseModel
	Follower string `gorm:"follower"`
	Followed string `gorm:"followed"`
}

func Follow(follower string, followed string) error {
	fs := FollowStruct{
		Follower: follower,
		Followed: followed,
	}
	err := db.Table("follows").Create(&fs).Error
	return err
}

func Unfollow(follower string, followed string) error {
	fs := FollowStruct{
		Follower: follower,
		Followed: followed,
	}
	fmt.Println(fs)
	err := db.Table("follows").Debug().Where(&fs).Delete(&FollowStruct{}).Error
	return err
}
