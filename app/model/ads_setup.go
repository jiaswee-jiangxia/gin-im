package model

import (
	"goskeleton/app/global/variable"
	"time"
)

type AdsSetupStruct struct {
	Id         int    `gorm:"column:id" json:"-"`
	Name       string `gorm:"column:name" json:"name"`
	ImgPreview string `gorm:"column:img_preview" json:"img_preview"`
	AdsType    string `gorm:"column:ads_type" json:"ads_type"`
	Content    string `gorm:"column:content" json:"content"`
}

func GetOnGoingAds() ([]*AdsSetupStruct, error) {
	var ads []*AdsSetupStruct
	now := time.Now().Format(variable.DateFormat)
	db := GetDB()
	query := db.Table("ads_setup as a").
		Where("a.start_at <= ?", now).
		Where("a.end_at >= ?", now)
	err := query.Find(&ads).Error
	if err != nil {
		return nil, err
	}
	return ads, nil
}
