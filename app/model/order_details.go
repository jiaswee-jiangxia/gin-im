package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type OrderDetailStruct struct {
	Id           int        `gorm:"column:id" json:"-"`
	OrderNo      string     `gorm:"column:order_no" json:"order_no"`
	MemberId     int        `gorm:"column:member_id" json:"member_id"`
	SerialNo     string     `gorm:"column:serial_no" json:"serial_no"`
	CartDetails  string     `gorm:"column:cart_details" json:"cart_details"`
	OrderType    int        `gorm:"column:order_type" json:"order_type"`
	RiderId      int        `gorm:"column:rider_id" json:"rider_id"`
	DeliveryAddr string     `gorm:"column:delivery_addr" json:"delivery_addr"`
	DeliveryLng  float64    `gorm:"column:delivery_lng" json:"delivery_lng"`
	DeliveryLat  float64    `gorm:"column:delivery_lat" json:"delivery_lat"`
	Status       int        `gorm:"column:status" json:"status"`
	CompletedAt  *time.Time `gorm:"column:completed_at" json:"completed_at"`
	CancelledAt  *time.Time `gorm:"column:cancelled_at" json:"cancelled_at"`
}

func CreateOrder(tx *gorm.DB, orderNo string, memId int, storeNo string, cartDetails interface{}, orderType int, riderId int, deliveryAddr string, deliveryLng float64, deliveryLat float64, status int) (*OrderDetailStruct, error) {
	text, _ := json.Marshal(cartDetails)
	order := OrderDetailStruct{
		OrderNo:      orderNo,
		MemberId:     memId,
		SerialNo:     storeNo,
		CartDetails:  string(text),
		OrderType:    orderType,
		RiderId:      riderId,
		DeliveryAddr: deliveryAddr,
		DeliveryLng:  deliveryLng,
		DeliveryLat:  deliveryLat,
		Status:       status,
		CompletedAt:  nil,
		CancelledAt:  nil,
	}
	err := tx.Table("order_details").Create(&order).Error

	if err != nil {
		return nil, err
	}
	return &order, nil
}

type OrderDetailListStruct struct {
	Id           int        `gorm:"column:id" json:"-"`
	OrderNo      string     `gorm:"column:order_no" json:"order_no"`
	MemberId     int        `gorm:"column:member_id" json:"member_id"`
	SerialNo     string     `gorm:"column:serial_no" json:"serial_no"`
	StoreName    string     `gorm:"column:store_name" json:"store_name"`
	StoreSubname string     `gorm:"column:store_subname" json:"store_subname"`
	StoreDesc    string     `gorm:"column:store_desc" json:"store_desc"`
	CartDetails  string     `gorm:"column:cart_details" json:"cart_details"`
	OrderType    int        `gorm:"column:order_type" json:"order_type"`
	RiderId      int        `gorm:"column:rider_id" json:"rider_id"`
	DeliveryAddr string     `gorm:"column:delivery_addr" json:"delivery_addr"`
	DeliveryLng  string     `gorm:"column:delivery_lng" json:"delivery_lng"`
	DeliveryLat  string     `gorm:"column:delivery_lat" json:"delivery_lat"`
	Status       string     `gorm:"column:status" json:"status"`
	CompletedAt  *time.Time `gorm:"column:completed_at" json:"completed_at"`
	CancelledAt  *time.Time `gorm:"column:cancelled_at" json:"cancelled_at"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
}

func GetOrderList(status string, page int, limit int, docNo string) ([]*OrderDetailListStruct, error) {
	var orders []*OrderDetailListStruct
	offset := (page - 1) * limit
	query := db.Table("order_details as a").
		Joins("inner join `db-cm-store`.store_detail b ON a.serial_no = b.serial_no")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if docNo != "" {
		query = query.Where("a.order_no = ?", docNo)
	}
	err := query.Limit(limit).Offset(offset).Select("a.*, b.store_name, b.store_subname, b.store_desc").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func UpdateCallBackOrderStatus(storeNo string, orderNo string, status string) error {
	err := db.Table("order_details as a").
		Where("a.serial_no = ?", storeNo).
		Where("a.order_no = ?", orderNo).
		Updates(map[string]interface{}{
			"status": status,
		}).Error

	if err != nil {
		return err
	}
	return nil
}
