package black_users_service

import (
	"gorm.io/gorm"
	"goskeleton/app/model"
	"strconv"
)

type BlackUsersStruct struct {
	UserId      string
	BlackUserId string
	Page        string
	Limit       string
}

func (m *BlackUsersStruct) CreateNewBlackUser() (*model.BlackUsers, error) {
	var u *model.BlackUsers
	userId, _ := strconv.Atoi(m.UserId)
	blackUserId, _ := strconv.Atoi(m.BlackUserId)
	_, err := model.GetBlackUser(m.UserId, m.BlackUserId)
	if err == gorm.ErrRecordNotFound {
		u, err = model.CreateNewBlackUser(&model.BlackUsers{
			UserId:      int64(userId),
			BlackUserId: int64(blackUserId),
			Active:      int64(1),
		})
		if err != nil {
			return nil, err
		}
	} else {
		u, err = model.BlackUsers{
			UserId:      int64(userId),
			BlackUserId: int64(blackUserId),
		}.Updates(map[string]interface{}{
			"active": 1,
		})
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (m *BlackUsersStruct) RemoveBlackUser() (*model.BlackUsers, error) {
	userId, _ := strconv.Atoi(m.UserId)
	blackUserId, _ := strconv.Atoi(m.BlackUserId)
	u, err := model.BlackUsers{
		UserId:      int64(userId),
		BlackUserId: int64(blackUserId),
	}.Updates(map[string]interface{}{
		"active": 0,
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *BlackUsersStruct) QueryBlackUser() ([]*model.BlackUsersUser, error) {
	page, _ := strconv.Atoi(m.Page)
	limit, _ := strconv.Atoi(m.Limit)
	limitStart := (page - 1) * limit
	userId, _ := strconv.Atoi(m.UserId)
	u, err := model.GetBlackUserList(&model.BlackUsers{UserId: int64(userId)}, limit, limitStart)
	if err != nil {
		return nil, err
	}
	return u, nil
}
