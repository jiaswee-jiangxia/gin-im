package model

import (
	"gorm.io/gorm"
	"time"
)

type StoreKycFiles struct {
	StoreId			int        `gorm:"column:store_id" json:"store_id"`
	FileType		string     `gorm:"column:file_type" json:"file_type"`
	FileUrl			string     `gorm:"column:file_url" json:"file_url"`
	Status			string     `gorm:"column:status" json:"status"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func CreateStoreKycFiles(tx *gorm.DB, storeId int, fileType string, fileUrl string) (*StoreKycFiles, error) {
	kycFiles := StoreKycFiles{
		StoreId: storeId,
		FileType: fileType,
		FileUrl: fileUrl,
		Status: "P",
		CreatedAt:  nil,
		UpdatedAt:  nil,
	}
	result := tx.Table("store_kyc_files").Where(StoreKycFiles{StoreId: storeId, FileType: fileType}).FirstOrCreate(&kycFiles)
	err := result.Error
	if err != nil {
		return nil, err
	}

	if result.RowsAffected == 0{
		kycFiles.FileUrl = fileUrl
		err = tx.Table("store_kyc_files").
			Where(StoreKycFiles{StoreId: storeId, FileType: fileType}).
			Updates(kycFiles).Error
		if err != nil {
			return nil, err
		}
	}
	return &kycFiles, nil
}