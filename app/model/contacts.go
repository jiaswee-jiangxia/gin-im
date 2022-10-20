package model

import (
	"encoding/json"
	"goskeleton/app/service/redis_service"
	"strconv"
)

// Contacts struct
type Contacts struct {
	BaseModel
	UserId   int64  `gorm:"column:user_id" json:"user_id"`
	FriendId int64  `gorm:"column:friend_id" json:"friend_id"`
	Status   int64  `gorm:"column:status" json:"status"`
	Grouping string `gorm:"column:grouping" json:"grouping"`
}

type UserContacts struct {
	Contacts
	BaseModel
	Username string `gorm:"column:username" json:"username"`
}

func GetContactsByBothId(userId string, friendId string) (*Contacts, error) {
	var u *Contacts
	var err error
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_CONTACT:" + userId + "-" + friendId,
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	cacheData := rdb.PrepareCacheRead()
	if cacheData != "" {
		err = json.Unmarshal([]byte(cacheData), &u)
	} else {
		err = db.Table("contacts").
			Where("user_id = ?", userId).
			Where("friend_id = ?", friendId).
			First(&u).Error
		if err == nil {
			rdb.CacheValue = u
			rdb.PrepareCacheWrite()
		}
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func CreateNewContact(u *Contacts) (*Contacts, error) {
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_CONTACT:" + strconv.Itoa(int(u.UserId)) + "-" + strconv.Itoa(int(u.FriendId)),
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	err := db.Table("contacts").
		Create(&u).Error
	if err != nil {
		return nil, err
	}
	rdb.CacheValue = u
	rdb.PrepareCacheWrite()
	return u, nil
}

func Updates(g *Contacts, updates interface{}) (*Contacts, error) {
	err := db.Table("contacts").Model(&g).
		Where("user_id = ?", g.UserId).
		Where("friend_id = ?", g.FriendId).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_CONTACT:" + strconv.Itoa(int(g.UserId)) + "-" + strconv.Itoa(int(g.FriendId)),
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	_ = rdb.DelCache()
	return g, nil
}

func GetContactList(u *Contacts) ([]*UserContacts, error) {
	var g []*UserContacts
	var err error
	rdb := redis_service.RedisStruct{
		CacheName:      "USER_CONTACT_LIST:" + strconv.Itoa(int(u.UserId)),
		CacheNameIndex: redis_service.RedisCacheUser,
	}
	cacheData := rdb.PrepareCacheRead()
	if cacheData != "" {
		err = json.Unmarshal([]byte(cacheData), &g)
		if err != nil {
			return nil, err
		}
	} else {
		err = db.Table("contacts").
			Joins("inner join users on users.id = contacts.friend_id").
			Where("user_id = ?", u.UserId).
			Where("status = ?", 1).
			Select("contacts.*, users.username").
			Find(&g).Error
	}
	if err != nil {
		return nil, err
	}
	rdb.CacheValue = g
	rdb.PrepareCacheWrite()
	return g, nil
}
