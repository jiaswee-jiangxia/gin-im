package group_service

import (
	"encoding/json"
	"goskeleton/app/model"
	redis "goskeleton/app/service/redis_service"
	"strconv"
)

type GroupStruct struct {
	model.BaseModel
	Name      string `gorm:"column:name" json:"groupname"`
	CreatedBy string `gorm:"column:created_by" json:"created_by"`
	Owner     string `gorm:"column:owner" json:"owner"`
	Disbanded bool   `gorm:"column:disbanded" json:"disbanded"`
}

type GroupMemberStruct struct {
	model.BaseModel
	GroupID  int64  `gorm:"column:group_id" json:"group_id"`
	Username string `gorm:"column:username" json:"user_id"`
	Role     string `gorm:"column:role" json:"role"`
}

func (m *GroupStruct) CreateGroup() (*model.GroupStruct, error) {
	group, err := model.CreateGroup(m.Name, m.Owner)
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

func (m *GroupStruct) GetGroupInfo() (g *model.GroupStruct, err error) {
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_INFO",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	cacheData := redisService.PrepareCacheRead()
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &g)
		return g, nil
	}
	g, err = model.GetGroupInfo(m.Id)
	if err != nil {
		return nil, err
	}
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return
}

func (m *GroupStruct) GetGroupAdmin() (g []model.GroupMemberStruct, err error) {
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_ADMIN",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	cacheData := redisService.PrepareCacheRead()
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &g)
		return g, nil
	}
	g, err = model.GetGroupAdminInfo(m.Id)
	if err != nil {
		return nil, err
	}
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return
}

func (m *GroupStruct) GetGroupMember() (g []model.GroupMemberStruct, err error) {
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_MEMBER",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	cacheData := redisService.PrepareCacheRead()
	if cacheData != "" {
		_ = json.Unmarshal([]byte(cacheData), &g)
		return g, nil
	}
	g, err = model.GetGroupMemberInfo(m.Id)
	if err != nil {
		return nil, err
	}
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return
}

func (m *GroupStruct) AddGroupMember(username string) error {
	err := model.AddGroupMember(m.Id, username)
	if err != nil {
		return err
	}
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_MEMBER",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	g, err := model.GetGroupMemberInfo(m.Id)
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return nil
}

func (m *GroupStruct) SetGroupAdmin(memberUsername string) error {
	err := model.SetGroupAdmin(m.Id, memberUsername)
	if err != nil {
		return err
	}
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_ADMIN",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	g, err := model.GetGroupAdminInfo(m.Id)
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return nil
}

func (m *GroupStruct) SetGroupOwner(memberUsername string) error {
	err := model.SetGroupOwner(m.Id, memberUsername)
	if err != nil {
		return err
	}
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_INFO",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	g, err := model.GetGroupInfo(m.Id)
	if err != nil {
		return err
	}
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return nil
}

func (m *GroupStruct) RemoveGroupMember(memberUsename string) error {
	err := model.RemoveGroupMember(m.Id, memberUsename)
	if err != nil {
		return err
	}
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_MEMBER",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	g, err := model.GetGroupMemberInfo(m.Id)
	if err != nil {
		return err
	}
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()
	return nil
}

func (m *GroupStruct) DisbandGroup() error {
	err := model.DisbandGroup(m.Id)
	if err != nil {
		return err
	}
	// Clearing Group info ---------------------------------------
	redisService := redis.RedisStruct{
		CacheName:      "GROUP_INFO",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	g, err := model.GetGroupInfo(m.Id)
	if err != nil {
		return err
	}
	redisService.CacheValue = g
	redisService.PrepareCacheWrite()

	// Clearing Group member info --------------------------------
	redisService = redis.RedisStruct{
		CacheName:      "GROUP_MEMBER",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	g2, err := model.GetGroupMemberInfo(m.Id)
	if err != nil {
		return err
	}
	redisService.CacheValue = g2
	redisService.PrepareCacheWrite()

	// Clearing Group admin info ---------------------------------
	redisService = redis.RedisStruct{
		CacheName:      "GROUP_ADMIN",
		CacheNameIndex: redis.RedisCacheGroup,
		CacheKey:       strconv.FormatInt(m.Id, 10),
	}
	g3, err := model.GetGroupAdminInfo(m.Id)
	if err != nil {
		return err
	}
	redisService.CacheValue = g3
	redisService.PrepareCacheWrite()
	return nil
}
