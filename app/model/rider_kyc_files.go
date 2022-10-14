package model

import (
	"gorm.io/gorm"
	"time"
)

type RiderKycFiles struct {
	RiderId			int        `gorm:"column:rider_id" json:"rider_id"`
	FileType		string     `gorm:"column:file_type" json:"file_type"`
	FileUrl			string     `gorm:"column:file_url" json:"file_url"`
	Status			string     `gorm:"column:status" json:"status"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func CreateRiderKycFiles(tx *gorm.DB, riderId int, fileType string, fileUrl string) (*RiderKycFiles, error) {
	kycFiles := RiderKycFiles{
		RiderId: riderId,
		FileType: fileType,
		FileUrl: fileUrl,
		Status: "P",
		CreatedAt:  nil,
		UpdatedAt:  nil,
	}
	result := tx.Table("rider_kyc_files").Where(RiderKycFiles{RiderId: riderId, FileType: fileType}).FirstOrCreate(&kycFiles)
	err := result.Error
	if err != nil {
		return nil, err
	}

	if result.RowsAffected == 0{
		kycFiles.FileUrl = fileUrl
		err = tx.Table("rider_kyc_files").
			Where(RiderKycFiles{RiderId: riderId, FileType: fileType}).
			Updates(kycFiles).Error
		if err != nil {
			return nil, err
		}
	}
	return &kycFiles, nil
}