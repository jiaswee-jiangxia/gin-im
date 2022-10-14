package redis_service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"log"
	"strings"
	"time"
)

const CacheNamePre = "im_"

const RedisCacheDefault = 0
const RedisCacheUser = 1

var RedisCacheExpired = map[string]time.Duration{
	// 永远不过期的redis为 -1
	"STORE_DETAILS":        -1,
	"STORE_MENU_DETAILS":   -1,
	"MEMBER_CART_CHECKOUT": 1800,
}

type RedisStruct struct {
	CacheName        string
	CacheNewName     string
	CacheValue       interface{}
	CacheNameIndex   int
	CacheKey         string
	OffJsonFormatter int
	CacheTimeInSec   int
}

func (m *RedisStruct) PrepareCacheRead() string {
	rdb := model.GetRedis(m.CacheNameIndex)
	ctx := context.TODO()
	var value *redis.StringCmd
	if m.CacheKey != "" {
		value = rdb.HGet(ctx, CacheNamePre+strings.ToLower(m.CacheName), strings.ToLower(m.CacheKey))
	} else {
		value = rdb.Get(ctx, CacheNamePre+strings.ToLower(m.CacheName))
	}
	return value.Val()
}

func (m *RedisStruct) PrepareCacheWrite() bool {
	rdb := model.GetRedis(m.CacheNameIndex)
	ctx := context.TODO()
	duration := 0 * time.Second
	cacheName := strings.Split(m.CacheName, ":")
	if len(cacheName) > 1 {
		if RedisCacheExpired[cacheName[0]] > 0 {
			duration = RedisCacheExpired[cacheName[0]] * time.Second
		} else if RedisCacheExpired[cacheName[0]] == 0 {
			if m.CacheTimeInSec < 1 {
				duration = -1
			} else {
				duration = time.Duration(m.CacheTimeInSec) * time.Second
			}
		}
	} else {
		if RedisCacheExpired[m.CacheName] > 0 {
			duration = RedisCacheExpired[cacheName[0]] * time.Second
		} else if RedisCacheExpired[m.CacheName] == 0 {
			if m.CacheTimeInSec < 1 {
				duration = -1
			} else {
				duration = time.Duration(m.CacheTimeInSec) * time.Second
			}
		}
	}

	var err error
	if m.CacheValue == "[]" || m.OffJsonFormatter == 1 {
		if m.CacheKey != "" {
			err = rdb.HSet(ctx, CacheNamePre+strings.ToLower(m.CacheName), strings.ToLower(m.CacheKey), m.CacheValue).Err()
			if err == nil {
				if duration > 0 {
					rdb.Expire(ctx, CacheNamePre+strings.ToLower(m.CacheName), duration)
				}
			}
		} else {
			err = rdb.Set(ctx, CacheNamePre+strings.ToLower(m.CacheName), m.CacheValue, duration).Err()
		}
	} else {
		content, errJson := json.Marshal(m.CacheValue)
		if errJson != nil {
			log.Fatalf("redis.PrepareCacheWrite err: json.Marshal")
		}
		if m.CacheKey != "" {
			err = rdb.HSet(ctx, CacheNamePre+strings.ToLower(m.CacheName), strings.ToLower(m.CacheKey), string(content)).Err()
			if err == nil {
				if duration > 0 {
					rdb.Expire(ctx, CacheNamePre+strings.ToLower(m.CacheName), duration)
				}
			}
		} else {
			err = rdb.Set(ctx, CacheNamePre+strings.ToLower(m.CacheName), string(content), duration).Err()
		}
	}
	if err != nil {
		variable.ZapLog.Error(err.Error())
		return false
	}
	return true
}

func (m *RedisStruct) PrepareCacheRename() string {
	rdb := model.GetRedis(m.CacheNameIndex)
	ctx := context.TODO()
	var value *redis.StatusCmd
	value = rdb.Rename(ctx, CacheNamePre+strings.ToLower(m.CacheName), CacheNamePre+strings.ToLower(m.CacheNewName))
	return value.Name()
}

func (m *RedisStruct) GetKeyRange() map[string]string {
	rdb := model.GetRedis(m.CacheNameIndex)
	ctx := context.TODO()
	value := rdb.HGetAll(ctx, CacheNamePre+strings.ToLower(m.CacheName))
	return value.Val()
}

func (m *RedisStruct) DelCache() error {
	var err error
	rdb := model.GetRedis(m.CacheNameIndex)
	ctx := context.TODO()
	if m.CacheKey != "" {
		err = rdb.HDel(ctx, CacheNamePre+strings.ToLower(m.CacheName), strings.ToLower(m.CacheKey)).Err()
	} else {
		err = rdb.Del(ctx, CacheNamePre+strings.ToLower(m.CacheName)).Err()
	}
	return err
}
