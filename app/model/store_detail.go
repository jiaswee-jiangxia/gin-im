package model

import (
	"gorm.io/gorm"
	"time"
)

type StoreDetailStruct struct {
	Id                 string  `gorm:"column:id" json:"-"`
	SerialNo           string  `gorm:"column:serial_no" json:"serial_no"`
	StoreName          string  `gorm:"column:store_name" json:"store_name"`
	StoreSubname       string  `gorm:"column:store_subname" json:"store_subname"`
	StoreDesc          string  `gorm:"column:store_desc" json:"store_desc"`
	Rate               string  `gorm:"column:rate" json:"rate"`
	ImageSet           string  `gorm:"column:img_set" json:"img_set"`
	Tags               string  `gorm:"column:tags" json:"tags"`
	Lat                float64 `gorm:"column:lat" json:"lat"`
	Lng                float64 `gorm:"column:lng" json:"lng"`
	OperationStatus    string  `gorm:"column:operation_status" json:"operation_status"`
	BasicPaid          float64 `gorm:"column:basic_paid" json:"basic_paid"`
	DistancePaid       float64 `gorm:"column:distance_paid" json:"distance_paid"`
	MenuCategorySorted string  `gorm:"column:menu_category_sorted" json:"menu_category_sorted"`
}

func GetStoreDetail(storeNo string) (*StoreDetailStruct, error) {
	var Store *StoreDetailStruct
	query := db.Table("stores as a").
		Joins("inner join store_detail b ON a.serial_no = b.serial_no").
		Where("b.serial_no = ?", storeNo)

	err := query.Order("scheme_id").Select("a.*, b.menu_category_sorted, b.store_name, b.store_subname, b.store_desc, b.basic_paid, b.distance_paid").First(&Store).Error
	if err != nil {
		return nil, err
	}
	return Store, nil
}

type StoreMenuStruct struct {
	Id           int     `gorm:"column:id" json:"id"`
	SerialNo     string  `gorm:"column:serial_no" json:"serial_no"`
	Category     string  `gorm:"column:category" json:"category"`
	CategoryName string  `gorm:"column:category_name" json:"category_name"`
	ProductName  string  `gorm:"column:product_name" json:"product_name"`
	ProductDesc  string  `gorm:"column:product_desc" json:"product_desc"`
	FileUrl      string  `gorm:"column:file_url" json:"file_url"`
	UnitPrice    float64 `gorm:"column:unit_price" json:"unit_price"`
	CheckQty     int     `gorm:"column:check_qty" json:"check_qty"`
	Qty          int     `gorm:"column:qty" json:"qty"`
	BChoices     int     `gorm:"column:b_choices" json:"b_choices"`
	Status       string  `gorm:"column:status" json:"status"`
}

type StoreFullMenuStruct struct {
	Id          int         `gorm:"column:id" json:"id"`
	SerialNo    string      `gorm:"column:serial_no" json:"-"`
	Category    string      `gorm:"column:category" json:"-"`
	ProductName string      `gorm:"column:product_name" json:"product_name"`
	ProductDesc string      `gorm:"column:product_desc" json:"product_desc"`
	FileUrl     string      `gorm:"column:file_url" json:"file_url"`
	UnitPrice   float64     `gorm:"column:unit_price" json:"unit_price"`
	CheckQty    int         `gorm:"column:check_qty" json:"check_qty"`
	Qty         int         `gorm:"column:qty" json:"qty"`
	BChoices    int         `gorm:"column:b_choices" json:"b_choices"`
	Status      string      `gorm:"column:status" json:"status"`
	Variety     interface{} `gorm:"column:variety" json:"variety"`
}

func GetMenuList(storeId int, storeNo string, categoryId string) ([]*StoreMenuStruct, error) {
	var Store []*StoreMenuStruct
	//index := storeId % 100
	query := db.Table("store_menu as a").
		Joins("inner join store_detail b ON a.serial_no = b.serial_no").
		Where("b.serial_no = ?", storeNo)

	if categoryId != "" {
		query = query.Where("category = ?", categoryId)
	}

	err := query.Order("a.id").Select("a.*").Find(&Store).Error
	if err != nil {
		return nil, err
	}
	return Store, nil
}

type StoreDetailStruct2 struct {
	Id            int        `gorm:"column:id" json:"-"`
	SerialNo      string     `gorm:"column:serial_no" json:"serial_no"`
	StoreName     string     `gorm:"column:store_name" json:"store_name"`
	StoreSubname  string     `gorm:"column:store_subname" json:"store_subname"`
	StoreDesc     string     `gorm:"column:store_desc" json:"store_desc"`
	StoreAddress  string     `gorm:"column:store_address" json:"store_address"`
	StoreCity     string     `gorm:"column:store_city" json:"store_city"`
	StoreState    string     `gorm:"column:store_state" json:"store_state"`
	StorePostcode string     `gorm:"column:store_postcode" json:"store_postcode"`
	StoreCountry  string     `gorm:"column:store_country" json:"store_country"`
	MapsRouting   string     `gorm:"column:maps_routing" json:"maps_routing"`
	BasicPaid     float64    `gorm:"column:basic_paid" json:"basic_paid"`
	DistancePaid  float64    `gorm:"column:distance_paid" json:"distance_paid"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func CreateStoreDetail(tx *gorm.DB, storeId int, serialNo string, storeName string, storeSubname string, storeDesc string, storeAddress string, storeCity string, storeState string, storePostcode string, storeCountry string) (*StoreDetailStruct2, error) {
	storeDetail := StoreDetailStruct2{
		Id:            storeId,
		SerialNo:      serialNo,
		StoreName:     storeName,
		StoreSubname:  storeSubname,
		StoreDesc:     storeDesc,
		StoreAddress:  storeAddress,
		StoreCity:     storeCity,
		StoreState:    storeState,
		StorePostcode: storePostcode,
		StoreCountry:  storeCountry,
		MapsRouting:   "",
		BasicPaid:     2,
		DistancePaid:  0.1,
		CreatedAt:     nil,
		UpdatedAt:     nil,
	}
	err := tx.Table("store_detail").Create(&storeDetail).Error

	if err != nil {
		return nil, err
	}
	return &storeDetail, nil
}
