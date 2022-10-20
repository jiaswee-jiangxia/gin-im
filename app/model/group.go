package model

import (
	"encoding/json"
	redis "goskeleton/app/service/redis_service"
	"strconv"
)

type GroupStruct struct {
	BaseModel
	Name      string `gorm:"column:name" json:"groupname"`
	CreatedBy string `gorm:"column:created_by" json:"created_by"`
	Owner     string `gorm:"column:owner" json:"owner"`
	Disbanded bool   `gorm:"column:disbanded" json:"disbanded"`
}

type GroupMemberStruct struct {
	BaseModel
	GroupID  int64  `gorm:"column:group_id" json:"group_id"`
	Username string `gorm:"column:username" json:"user_id"`
	Role     string `gorm:"column:role" json:"role"`
}

func CreateGroup(groupName string, owner string) (*GroupStruct, error) {
	group := &GroupStruct{
		Name:      groupName,
		CreatedBy: owner,
		Owner:     owner,
		Disbanded: false,
	}
	err := db.Table("groups").Create(&group).Error
	if err != nil {
		return nil, err
	}
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_INFO",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(group.Id, 10),
	}
	redisService.CacheValue = group
	redisService.PrepareCacheWrite()
	return group, nil
}

func AddGroupMember(groupID int64, memberUsername string) error {
	err := db.Table("group_members").Select("group_id", "username", "role").
		Create(&GroupMemberStruct{
			GroupID:  groupID,
			Username: memberUsername,
			Role:     "member",
		}).Error
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_MEMBER",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	redisService.DelCache()
	return err
}

func GetGroupInfo(groupID int64) (g *GroupStruct, err error) {
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_INFO",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	cacheData := redisService.PrepareCacheRead()
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &g)
		return g, nil
	}
	g = &GroupStruct{}
	err = db.Table("groups").Where("id", groupID).Scan(&g).Error
	if err != nil {
		return nil, err
	}
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return
}

func GetGroupAdminInfo(groupID int64) (g []GroupMemberStruct, err error) {
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_ADMIN",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	cacheData := redisService.PrepareCacheRead()
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &g)
		return g, nil
	}
	g = make([]GroupMemberStruct, 0)
	err = db.Table("group_members").Where("group_id", groupID).Where("role", "admin").Scan(&g).Error
	if err != nil {
		return nil, err
	}
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return
}

func GetGroupMemberInfo(groupID int64) (g []GroupMemberStruct, err error) {
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_MEMBER",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	cacheData := redisService.PrepareCacheRead()
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &g)
		return g, nil
	}
	g = make([]GroupMemberStruct, 0)
	err = db.Table("group_members").Where("group_id", groupID).Scan(&g).Error
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return
}

func SetGroupAdmin(groupID int64, memberUsername string) error {
	err := db.Table("group_members").Where("group_id", groupID).Where("username", memberUsername).Update("role", "admin").Error
	if err != nil {
		return err
	}
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_ADMIN",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	redisService.DelCache()
	return err
}

func SetGroupOwner(groupID int64, memberUsername string) error {
	err := db.Table("groups").Where("id", groupID).Update("owner", memberUsername).Error
	if err != nil {
		return err
	}
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_INFO",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	redisService.DelCache()
	return err
}

func RemoveGroupMember(groupID int64, memberUsername string) error {
	member := &GroupMemberStruct{}
	err := db.Table("group_members").Where("group_id", groupID).Where("username", memberUsername).Delete(&member).Error
	if err != nil {
		return err
	}
	// Delete member cache
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_MEMBER",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	redisService.DelCache()

	// Delete admin cache
	redisService = redis.RedisStruct{
		CacheName:      "GROUP_ADMIN",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	redisService.DelCache()
	return err
}

func DisbandGroup(groupID int64) error {
	member := &GroupMemberStruct{}
	err := db.Table("group_members").Where("group_id", groupID).Delete(&member).Error
	if err != nil {
		return err
	}
	err = db.Table("groups").Where("id", groupID).Update("disbanded", true).Error
	if err != nil {
		return err
	}
	// Clearing Group info ---------------------------------------
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_INFO",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	redisService.DelCache()

	// Clearing Group member info --------------------------------
	redisService = redis.RedisStruct{
		CacheName:      "GROUP_MEMBER",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	redisService.DelCache()

	// Clearing Group admin info ---------------------------------
	redisService = redis.RedisStruct{
		CacheName:      "GROUP_ADMIN",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(groupID, 10),
	}
	redisService.DelCache()
	return err
}
