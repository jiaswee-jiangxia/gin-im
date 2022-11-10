package follow_service

import "goskeleton/app/model"

type FollowStruct struct {
	Follower string `json:"follower"`
	Followed string `json:"followed"`
}

func (m *FollowStruct) Follow() error {
	err := model.Follow(m.Follower, m.Followed)
	return err
}

func (m *FollowStruct) Unfollow() error {
	err := model.Unfollow(m.Follower, m.Followed)
	return err
}

func (m *FollowStruct) GetMyFollowList() ([]string, error) {
	followList, err := model.GetMyFollowList(m.Follower)
	return followList, err
}

func (m *FollowStruct) GetUserFollowCount() int {
	count := model.GetUserFollowCount(m.Follower)
	return count
}
