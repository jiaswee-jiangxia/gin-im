package model

import (
	"gorm.io/gorm"
	"time"
)

type MerchantKycFiles struct {
	MerchantId			int        `gorm:"column:merchant_id" json:"merchant_id"`
	FileType		string     `gorm:"column:file_type" json:"file_type"`
	FileUrl			string     `gorm:"column:file_url" json:"file_url"`
	Status			string     `gorm:"column:status" json:"status"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func CreateMerchantKycFiles(tx *gorm.DB, merchantId int, fileType string, fileUrl string) (*MerchantKycFiles, error) {
	kycFiles := MerchantKycFiles{
		MerchantId: merchantId,
		FileType: fileType,
		FileUrl: fileUrl,
		Status: "P",
		CreatedAt:  nil,
		UpdatedAt:  nil,
	}
	result := tx.Table("merchant_kyc_files").Where(MerchantKycFiles{MerchantId: merchantId, FileType: fileType}).FirstOrCreate(&kycFiles)
	err := result.Error
	if err != nil {
		return nil, err
	}

	if result.RowsAffected == 0{
		kycFiles.FileUrl = fileUrl
		err = tx.Table("merchant_kyc_files").
			Where(MerchantKycFiles{MerchantId: merchantId, FileType: fileType}).
			Updates(kycFiles).Error
		if err != nil {
			return nil, err
		}
	}
	return &kycFiles, nil
}