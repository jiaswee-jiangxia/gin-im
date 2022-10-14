package model

import (
	"gorm.io/gorm"
	"goskeleton/app/helpers"
	"time"
)

type MerchantWalletDetailsStruct struct {
	Id           int       `gorm:"column:id" json:"-"`
	MemberId     string    `gorm:"column:member_id" json:"member_id"`
	DocNo        string    `gorm:"column:doc_no" json:"doc_no"`
	WalletTypeId string    `gorm:"column:wallet_type_id" json:"-"`
	TransType    string    `gorm:"column:trans_type" json:"trans_type"`
	TotalIn      float64   `gorm:"column:total_in" json:"total_in"`
	TotalOut     float64   `gorm:"column:total_out" json:"total_out"`
	LastBalance  float64   `gorm:"column:last_balance" json:"last_balance"`
	Status       string    `gorm:"column:status" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func CreateMerchantWalletDetail(tx *gorm.DB, memId string, walletTypeId string, transType string, amount float64, lastBal float64) (*MerchantWalletDetailsStruct, error) {
	docNo := helpers.GenDocNo(6, "")
	totalIn := 0.0
	totalOut := 0.0
	if amount > 0 {
		totalIn = amount
		lastBal += amount
	} else {
		totalOut = amount * -1
		lastBal -= totalOut
	}
	memWallet := MerchantWalletDetailsStruct{
		MemberId:     memId,
		DocNo:        docNo,
		WalletTypeId: walletTypeId,
		TransType:    transType,
		TotalIn:      totalIn,
		TotalOut:     totalOut,
		LastBalance:  lastBal,
		Status:       "A",
		CreatedAt:    time.Now(),
	}
	err := tx.Table("merchant_wallet_details").Create(&memWallet).Error
	if err != nil {
		return nil, err
	}
	return &memWallet, nil
}

func GetMerchantWalletStatement(memId string, page int, limit int) ([]*MerchantWalletDetailsStruct, error) {
	var ewtStatement []*MerchantWalletDetailsStruct
	offset := (page - 1) * limit
	query := db.Table("merchant_wallet_details as a").
		Where("member_id = ?", memId)

	err := query.Limit(limit).Offset(offset).Select("a.*").Find(&ewtStatement).Error
	if err != nil {
		return nil, err
	}
	return ewtStatement, nil
}
