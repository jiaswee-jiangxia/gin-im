package contacts_service

import (
	"goskeleton/app/model"
	"goskeleton/app/service/redis_service"
	"strconv"
)

type ContactsStruct struct {
	UserId   string
	FriendId string
	Status   int64
	Grouping string
}

func (m *ContactsStruct) GetContactsByBothId() (*model.Contacts, error) {
	u, err := model.GetContactsByBothId(m.UserId, m.FriendId)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *ContactsStruct) CreateNewContact() (*model.Contacts, error) {
	userId, _ := strconv.Atoi(m.UserId)
	frdId, _ := strconv.Atoi(m.FriendId)
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_CONTACT:" + strconv.Itoa(userId) + "-" + strconv.Itoa(frdId),
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	u, err := model.CreateNewContact(&model.Contacts{
		UserId:   int64(userId),
		FriendId: int64(frdId),
		Status:   m.Status,
		Grouping: m.Grouping,
	})
	if err != nil {
		return nil, err
	}
	rdb.CacheValue = u
	rdb.PrepareCacheWrite()
	if m.Status > 0 {
		rdb.CacheName = "USER_CONTACT:" + strconv.Itoa(frdId) + "-" + strconv.Itoa(userId)
		u2, err := model.CreateNewContact(&model.Contacts{
			UserId:   int64(frdId),
			FriendId: int64(userId),
			Status:   m.Status,
		})
		if err != nil {
			return nil, err
		}
		rdb.CacheValue = u2
		rdb.PrepareCacheWrite()
	}
	return u, nil
}

func (m *ContactsStruct) AcceptContact() (*model.Contacts, error) {
	userId, _ := strconv.Atoi(m.UserId)
	frdId, _ := strconv.Atoi(m.FriendId)
	u, err := model.Updates(&model.Contacts{
		UserId:   int64(userId),
		FriendId: int64(frdId),
	}, map[string]interface{}{
		"status": m.Status,
	})
	if err != nil {
		return nil, err
	}
	_, err = model.CreateNewContact(&model.Contacts{
		UserId:   int64(frdId),
		FriendId: int64(userId),
		Status:   m.Status,
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *ContactsStruct) UpdateContact() (*model.Contacts, error) {
	userId, _ := strconv.Atoi(m.UserId)
	frdId, _ := strconv.Atoi(m.FriendId)
	u, err := model.Updates(&model.Contacts{
		UserId:   int64(userId),
		FriendId: int64(frdId),
	}, map[string]interface{}{
		"status": m.Status,
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *ContactsStruct) UpdateContactGrouping() (*model.Contacts, error) {
	userId, _ := strconv.Atoi(m.UserId)
	frdId, _ := strconv.Atoi(m.FriendId)
	u, err := model.Updates(&model.Contacts{
		UserId:   int64(userId),
		FriendId: int64(frdId),
	}, map[string]interface{}{
		"grouping": m.Grouping,
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *ContactsStruct) GetContactList() ([]*model.UserContacts, error) {
	userId, _ := strconv.Atoi(m.UserId)
	frdId, _ := strconv.Atoi(m.FriendId)
	u, err := model.GetContactList(&model.Contacts{
		UserId:   int64(userId),
		FriendId: int64(frdId),
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}
