package model

import (
	"gorm.io/gorm"
)

type MerchantWalletSummaryStruct struct {
	//Id         int    `gorm:"column:id" json:"-"`
	MemberId     string  `gorm:"column:member_id" json:"member_id"`
	WalletTypeId string  `gorm:"column:wallet_type_id" json:"-"`
	TotalIn      float64 `gorm:"column:total_in" json:"total_in"`
	TotalOut     float64 `gorm:"column:total_out" json:"total_out"`
	TotalBalance float64 `gorm:"column:total_balance" json:"total_balance"`
	CheckSum     string  `gorm:"column:check_sum" json:"-"`
}

func GetMerchantMemberBalance(memId string, walletTypeId string) (*MerchantWalletSummaryStruct, error) {
	var bal *MerchantWalletSummaryStruct
	//now := time.Now().Format(variable.DateFormat)
	query := db.Table("merchant_wallet_summary as a").
		Where("a.member_id = ?", memId).
		Where("a.wallet_type_id = ?", walletTypeId)
	err := query.Find(&bal).Error
	if err != nil {
		return nil, err
	}
	return bal, nil
}

func MerchantWalletAdjustment(tx *gorm.DB, memId string, walletTypeId string, amount float64) bool {
	var query *gorm.DB
	//now := time.Now().Format(variable.DateFormat)
	//db := GetDB()
	if amount > 0 {
		query = tx.Table("merchant_wallet_summary as a").
			Where("a.member_id = ?", memId).
			Where("a.wallet_type_id = ?", walletTypeId).
			Updates(map[string]interface{}{
				"total_in":      gorm.Expr("total_in + ?", amount),
				"total_balance": gorm.Expr("total_balance + ?", amount),
			})
	} else if amount < 0 {
		amount *= -1
		query = tx.Table("merchant_wallet_summary as a").
			Where("a.member_id = ?", memId).
			Where("a.wallet_type_id = ?", walletTypeId).
			Updates(map[string]interface{}{
				"total_out":     gorm.Expr("total_out + ?", amount),
				"total_balance": gorm.Expr("total_balance - ?", amount),
			})
	}

	if query.RowsAffected > 0 {
		return true
	}
	return false
}

func CreateMerchantMemberWallet(tx *gorm.DB, memId string, walletTypeId string) (*MerchantWalletSummaryStruct, error) {
	memWallet := MerchantWalletSummaryStruct{
		MemberId:     memId,
		WalletTypeId: walletTypeId,
		TotalIn:      0,
		TotalOut:     0,
		TotalBalance: 0,
		CheckSum:     "123",
	}
	err := tx.Table("merchant_wallet_summary").Create(&memWallet).Error
	if err != nil {
		return nil, err
	}
	return &memWallet, nil
}
