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
	err := db.Table("follows").Where(&fs).Delete(&FollowStruct{}).Error
	return err
}

func GetMyFollowList(username string) ([]string, error) {
	var followList []FollowStruct
	var outList []string
	fs := FollowStruct{
		Follower: username,
	}
	err := db.Table("follows").Where(&fs).Find(&followList).Error
	for _, v := range followList {
		outList = append(outList, v.Followed)
	}
	return outList, err
}

func GetUserFollowCount(username string) int {
	var count int64
	fs := FollowStruct{
		Follower: username,
	}
	db.Table("follows").Where(&fs).Count(&count)
	return int(count)
}
