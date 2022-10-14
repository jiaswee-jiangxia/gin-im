package model

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"goskeleton/app/helpers"
	"time"
)

type HomeListStruct struct {
	Id              string  `gorm:"column:id" json:"-"`
	SerialNo        string  `gorm:"column:serial_no" json:"serial_no"`
	StoreName       string  `gorm:"column:store_name" json:"store_name"`
	StoreSubname    string  `gorm:"column:store_subname" json:"store_subname"`
	StoreDesc       string  `gorm:"column:store_desc" json:"store_desc"`
	Rate            string  `gorm:"column:rate" json:"rate"`
	ImageSet        string  `gorm:"column:img_set" json:"img_set"`
	Category        string  `gorm:"column:category" json:"category"`
	Tags            string  `gorm:"column:tags" json:"tags"`
	Lat             float64 `gorm:"column:lat" json:"lat"`
	Lng             float64 `gorm:"column:lng" json:"lng"`
	OperationStatus string  `gorm:"column:operation_status" json:"operation_status"`
	BasicPaid       float64 `gorm:"column:basic_paid" json:"basic_paid"`
	DistancePaid    float64 `gorm:"column:distance_paid" json:"distance_paid"`
}

type StoreKycStruct struct {
	Id         int    `gorm:"column:id" json:"-"`
	SerialNo   string `gorm:"column:serial_no" json:"serial_no"`
	MerchantId int    `gorm:"column:merchant_id" json:"merchant_id"`
	BKyc       int    `gorm:"column:b_kyc" json:"b_kyc"`
}

func GetHomeList(Lat string, Lng string, AreaType string, AreaName string, schemeId int, limit int) ([]*HomeListStruct, error) {
	var Stores []*HomeListStruct
	var AreaKey string
	if AreaType == "city" {
		AreaKey = "store_city"
	} else if AreaType == "state" {
		AreaKey = "store_state"
	} else {
		AreaKey = "store_country"
	}
	positionLat, errLat := helpers.ParseDecimal(Lat, 1)
	positionLng, errLng := helpers.ParseDecimal(Lng, 1)
	if errLat != nil || errLng != nil {
		return nil, errors.New("position_undefined")
	}
	query := db.Table("stores as a").
		Joins("inner join store_detail b ON a.serial_no = b.serial_no").
		Where("lat >= ?", positionLat).
		Where("lng <= ?", positionLng).
		Where("scheme_id = ?", schemeId)

	if AreaName != "" {
		query = query.Or(AreaKey+" = ?", AreaName).
			Where("scheme_id = ?", schemeId)
	}
	err := query.Limit(limit).Order("scheme_id").Select("a.*, b.store_name, b.store_subname, b.store_desc, b.basic_paid, b.distance_paid").Find(&Stores).Error
	if err != nil {
		return nil, err
	}
	return Stores, nil
}

func GetStoreList(Lat string, Lng string, AreaType string, AreaName string, CategoryId string, Keyword string, page int, limit int) ([]*HomeListStruct, error) {
	var Stores []*HomeListStruct
	var AreaKey string
	if AreaType == "city" {
		AreaKey = "store_city"
	} else if AreaType == "state" {
		AreaKey = "store_state"
	} else {
		AreaKey = "store_country"
	}
	positionLat, errLat := helpers.ParseDecimal(Lat, 1)
	positionLng, errLng := helpers.ParseDecimal(Lng, 1)
	if errLat != nil || errLng != nil {
		return nil, errors.New("position_undefined")
	}
	offset := (page - 1) * limit
	query := db.Table("stores as a").
		Joins("inner join store_detail b ON a.serial_no = b.serial_no").
		Where("lat >= ?", positionLat).
		Where("lng <= ?", positionLng)
	if CategoryId != "0" {
		query = query.Where("category LIKE ?", "%,"+CategoryId+",%")
	}
	if Keyword != "" {
		query = query.Where("b.store_name LIKE ?", "%"+Keyword+"%")
	}
	if AreaName != "" {
		query = query.Or(AreaKey+" = ?", AreaName)
		if CategoryId != "0" {
			query = query.Where("category LIKE ?", "%,"+CategoryId+",%")
		}
		if Keyword != "" {
			query = query.Where("b.store_name LIKE ?", "%"+Keyword+"%")
		}
	}
	err := query.Limit(limit).Offset(offset).Select("a.*, b.store_name, b.store_subname, b.store_desc, b.basic_paid, b.distance_paid").Find(&Stores).Error
	if err != nil {
		return nil, err
	}
	return Stores, nil
}

type StoreSearchStruct struct {
	StoreName string `gorm:"column:store_name" json:"store_name"`
}

func GetStoreSearch(Keyword string) ([]*StoreSearchStruct, error) {
	var Stores []*StoreSearchStruct
	query := db.Table("store_detail as a")
	if Keyword != "" {
		query = query.Where("a.store_name LIKE ?", "%"+Keyword+"%")
	}
	err := query.Limit(10).Select("a.store_name").Find(&Stores).Error
	if err != nil {
		return nil, err
	}
	return Stores, nil
}

func GetStoresByMerchantId(merchantId int, storeNo string, page int, limit int) ([]*HomeListStruct, error) {
	var Stores []*HomeListStruct
	offset := (page - 1) * limit
	query := db.Table("stores as a").
		Joins("inner join store_detail b ON a.serial_no = b.serial_no").
		Where("a.merchant_id = ?", merchantId)

	if storeNo != "" {
		query = query.Where("a.serial_no = ?", storeNo)
	}

	err := query.Limit(limit).Offset(offset).Select("a.*, b.store_name, b.store_subname, b.store_desc, b.basic_paid, b.distance_paid").Find(&Stores).Error
	if err != nil {
		return nil, err
	}
	return Stores, nil
}

func GetStoreByMerchantId(merchantId int, storeNo string) (*HomeListStruct, error) {
	var Stores *HomeListStruct
	query := db.Table("stores as a").
		Joins("inner join store_detail b ON a.serial_no = b.serial_no").
		Where("a.merchant_id = ?", merchantId)

	if storeNo != "" {
		query = query.Where("a.serial_no = ?", storeNo)
	}

	err := query.Select("a.*, b.store_name, b.store_subname, b.store_desc, b.basic_paid, b.distance_paid").First(&Stores).Error
	if err != nil {
		return nil, err
	}
	return Stores, nil
}

func UpdateStoreImages(merchantId int, storeId int, fileUrl string) error {
	jsonText, _ := json.Marshal([]string{fileUrl})
	flag := db.Table("stores as a").
		Where("a.merchant_id = ?", merchantId).
		Where("a.id = ?", storeId).
		Updates(map[string]interface{}{
			"img_set": string(jsonText),
		}).Error
	return flag
}

func UpdateStoreKycStatus(tx *gorm.DB, storeId int, status int) error {
	err := tx.Table("stores").
		Where(StoreKycStruct{Id: storeId}).
		Updates(map[string]interface{}{
			"b_kyc": status,
		}).Error
	return err
}

type StoreStruct struct {
	Id              int        `gorm:"column:id" json:"-"`
	SerialNo        string     `gorm:"column:serial_no" json:"serial_no"`
	MerchantId      int        `gorm:"column:merchant_id" json:"merchant_id"`
	Status          string     `gorm:"column:status" json:"status"`
	StartDate       *time.Time `gorm:"column:start_date" json:"start_date"`
	EndDate         *time.Time `gorm:"column:end_date" json:"end_date"`
	Rate            float64    `gorm:"column:rate" json:"rate"`
	ImgSet          string     `gorm:"column:img_set" json:"img_set"`
	Category        string     `gorm:"column:category" json:"category"`
	Tags            string     `gorm:"column:tags" json:"tags"`
	SchemeId        int        `gorm:"column:scheme_id" json:"scheme_id"`
	Lat             float64    `gorm:"column:lat" json:"lat"`
	Lng             float64    `gorm:"column:lng" json:"lng"`
	Badge           string     `gorm:"column:badge" json:"badge"`
	OperationStatus string     `gorm:"column:operation_status" json:"operation_status"`
	BKyc            int        `gorm:"column:b_kyc" json:"b_kyc"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func CreateStores(tx *gorm.DB, serialNo string, merchantId int, categories string, lat float64, lng float64) (*StoreStruct, error) {
	store := StoreStruct{
		SerialNo:        serialNo,
		MerchantId:      merchantId,
		Status:          "P",
		StartDate:       nil,
		EndDate:         nil,
		Rate:            0,
		ImgSet:          "[]",
		Category:        categories,
		Tags:            "[]",
		SchemeId:        0,
		Lat:             lat,
		Lng:             lng,
		Badge:           "[]",
		OperationStatus: "I",
		BKyc:            0,
		CreatedAt:       nil,
		UpdatedAt:       nil,
	}
	err := tx.Table("stores").Create(&store).Error

	if err != nil {
		return nil, err
	}
	return &store, nil
}
