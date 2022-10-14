package model

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"os"
)

var db *gorm.DB
var rdb [16]*redis.Client
var rdbCount = 16

func Setup() {
	db = variable.GormDbMysql
}

func SetupRedis() {
	sum := 0
	for i := 0; i < rdbCount; i++ {
		rdb[i] = ConnectRedis(i)
		sum += i
	}
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

func GetDB() *gorm.DB {
	return db
}

func GetRedis(index int) *redis.Client {
	return rdb[index]
}

// SQLDataPaginateStdReturn. use in standard return for sql pagination
type SQLPaginateStdReturn struct {
	CurrentPage           int64   `json:"current_page"`
	PerPage               int64   `json:"per_page"`
	TotalCurrentPageItems int64   `json:"total_current_page_items"`
	TotalPage             float64 `json:"total_page"`
	TotalPageItems        int64   `json:"total_page_items"`
}

func SaveTx(tx *gorm.DB, value interface{}) error {
	err := tx.Save(value).Error
	if err != nil {
		return err
	}
	return nil
}
