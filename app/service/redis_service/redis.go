package redis_service

import (
	"context"
	"encoding/json"
	"fmt"
	"goskeleton/app/global/variable"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const CacheNamePre = "im_"

const RedisCacheDefault = 0
const RedisCacheUser = 1
const RedisCacheGroup = 2
const RedisCacheLock = 3
const RedisCacheFollow = 4

var rdb [16]*redis.Client
var rdbCount = 16

var RedisCacheExpired = map[string]time.Duration{
	// 永远不过期的redis为 -1
	//"STORE_DETAILS":        -1,
	"USER_PROFILE": 86400, //24hrs
	"USER_CONTACT": 86400, //24hrs
	"GROUP_INFO":   86400, //24hrs
	"GROUP_ADMIN":  86400, //24hrs
	"GROUP_MEMBER": 86400, //24hrs
}

type RedisStruct struct {
	CacheName          string
	CacheNewName       string
	CacheValue         interface{}
	CacheNameIndex     int
	CacheKey           string
	OffJsonFormatter   int
	CacheTimeInSec     int
	CacheLockTimeInSec int
}

func SetupRedis() {
	sum := 0
	for i := 0; i < rdbCount; i++ {
		rdb[i] = ConnectRedis(i)
		sum += i
	}
}

func GetRedis(index int) *redis.Client {
	return rdb[index]
}

func ConnectRedis(index int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Username: os.Getenv("REDIS_USERNAME"),
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       index,
		//TLSConfig: &tls.Config{
		//	InsecureSkipVerify: true,
		//},
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		variable.ZapLog.Fatal(fmt.Sprintf("Cannot Ping: %v\n", err.Error()))
	}
	return client
}

func (m *RedisStruct) PrepareCacheRead() string {
	rdb := GetRedis(m.CacheNameIndex)
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
	rdb := GetRedis(m.CacheNameIndex)
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

func PrepareLockTrial(cacheNameIndex int, cacheName string, cacheLockValue interface{}, cacheLockTimeInSec int) bool {
	rdb := GetRedis(cacheNameIndex)
	ctx := context.TODO()
	tJson, err := json.Marshal(cacheLockValue)
	if err != nil {
		variable.ZapLog.Error(err.Error())
		return false
	}
	flag := rdb.SetNX(ctx, CacheNamePre+strings.ToLower(cacheName), string(tJson), time.Duration(cacheLockTimeInSec)*time.Second)
	boolFlag, err := flag.Result()
	if err != nil {
		variable.ZapLog.Error(err.Error())
		return false
	}
	return boolFlag
}

func (m *RedisStruct) PrepareCacheRename() string {
	rdb := GetRedis(m.CacheNameIndex)
	ctx := context.TODO()
	var value *redis.StatusCmd
	value = rdb.Rename(ctx, CacheNamePre+strings.ToLower(m.CacheName), CacheNamePre+strings.ToLower(m.CacheNewName))
	return value.Name()
}

func (m *RedisStruct) GetKeyRange() map[string]string {
	rdb := GetRedis(m.CacheNameIndex)
	ctx := context.TODO()
	value := rdb.HGetAll(ctx, CacheNamePre+strings.ToLower(m.CacheName))
	return value.Val()
}

func (m *RedisStruct) DelCache() error {
	var err error
	rdb := GetRedis(m.CacheNameIndex)
	ctx := context.TODO()
	if m.CacheKey != "" {
		err = rdb.HDel(ctx, CacheNamePre+strings.ToLower(m.CacheName), strings.ToLower(m.CacheKey)).Err()
	} else {
		err = rdb.Del(ctx, CacheNamePre+strings.ToLower(m.CacheName)).Err()
	}
	return err
}

func PrepareUnlockTrial(cacheNameIndex int, cacheName string) bool {
	rdb := GetRedis(cacheNameIndex)
	ctx := context.TODO()
	err := rdb.Del(ctx, CacheNamePre+strings.ToLower(cacheName)).Err()
	if err != nil {
		return false
	}
	return true
}
