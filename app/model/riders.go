package model

import (
	"fmt"
	"gorm.io/gorm"
)

type RiderStruct struct {
	Id        int    `gorm:"column:id" json:"-"`
	Username  string `gorm:"column:username" json:"username"`
	Email     string `gorm:"column:email" json:"email"`
	MobileNo  string `gorm:"column:mobile_no" json:"mobile_no"`
	BKyc      int    `gorm:"column:b_kyc" json:"b_kyc"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

func GetRiderProfile(userId string, username string) (*UserStruct, error) {
	var Member *UserStruct
	//idx := helpers.ShardHash(username)
	err := db.Table("riders").
		Where("id = ?", userId).
		Where("username = ?", username).
		First(&Member).Error
	if err != nil {
		return nil, err
	}
	return Member, nil
}

func UpdateRiderKycStatus(tx *gorm.DB, riderId int, status int) error {
	err := tx.Table("riders").
		Where(RiderStruct{Id: riderId}).
		Updates(map[string]interface{}{
			"b_kyc": status,
		}).Error

	fmt.Println(err, riderId, status)
	return err
}
