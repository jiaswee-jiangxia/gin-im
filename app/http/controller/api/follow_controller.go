package api

import (
	"fmt"
	consts "goskeleton/app/global/response"
	"goskeleton/app/service/follow_service"
	"goskeleton/app/utils/response"

	"github.com/gin-gonic/gin"
)

type TargetUser struct {
	Target string `json:"target"`
}

type Follow struct {
	Follower string `json:"follower"`
	Followed string `json:"followed"`
}

func FollowUser(context *gin.Context) {
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)
	var target TargetUser
	if err := context.ShouldBindJSON(&target); err != nil {
		response.ErrorParam(context, target)
		return
	}
	followService := &follow_service.FollowStruct{
		Follower: usernameText,
		Followed: target.Target,
	}
	err := followService.Follow()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	response.Success(context, consts.Success, nil)
	return
}

func UnfollowUser(context *gin.Context) {
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)
	var target TargetUser
	if err := context.ShouldBindJSON(&target); err != nil {
		response.ErrorParam(context, target)
		return
	}
	followService := &follow_service.FollowStruct{
		Follower: usernameText,
		Followed: target.Target,
	}
	err := followService.Unfollow()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	response.Success(context, consts.Success, nil)
	return
}

func GetMyFollowList(context *gin.Context) {
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)
	followService := &follow_service.FollowStruct{
		Follower: usernameText,
	}
	followList, err := followService.GetMyFollowList()
	if err != nil {
		response.SuccessButFail(context, err.Error(), consts.Failed, nil)
		return
	}
	response.Success(context, consts.Success, followList)
}

func GetUserFollowCount(context *gin.Context) {
	var target TargetUser
	var count int
	if err := context.ShouldBindJSON(&target); err != nil {
		response.ErrorParam(context, target)
		return
	}
	followService := &follow_service.FollowStruct{
		Follower: target.Target,
	}
	count = followService.GetUserFollowCount()
	response.Success(context, consts.Success, count)
	return
}

func CheckIsFriend(context *gin.Context) {
	username, exist := context.Get("username")
	if !exist {
		response.Success(context, consts.Failed, nil)
	}
	usernameText := fmt.Sprintf("%v", username)
	var target TargetUser
	if err := context.ShouldBindJSON(&target); err != nil {
		response.ErrorParam(context, target)
		return
	}

	check1 := CheckFollowed(usernameText, target.Target)
	check2 := CheckFollowed(target.Target, usernameText)
	var out bool = check1 && check2
	response.Success(context, consts.Success, out)
	return
}

func CheckFollowed(follower string, followed string) bool {
	// rdb := redis_service.RedisStruct{
	// 	CacheName:      "USER_FOLLOWED:" + follower,
	// 	CacheNameIndex: redis_service.RedisCacheFollow,
	// }
	// cacheData := rdb.PrepareCacheRead()
	// if cacheData != "" {
	// 	var returnList interface{}
	// 	_ = json.Unmarshal([]byte(cacheData), &returnList)
	// 	return true
	// }
	followService := &follow_service.FollowStruct{
		Follower: follower,
	}
	followList, err := followService.GetMyFollowList()
	if err != nil {
		return false
	}
	for _, v := range followList {
		if v == followed {
			return true
		}
	}
	return false
}
