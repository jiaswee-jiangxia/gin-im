package model

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
)

var db *gorm.DB

func Setup() {
	db = variable.GormDbMysql
}

func GetDB() *gorm.DB {
	return db
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
