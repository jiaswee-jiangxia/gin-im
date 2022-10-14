package model

type PaymentGatewaySetup struct {
	Id			int 	`gorm:"column:id" json:"-"`
	Name		string 	`gorm:"column:name" json:"name"`
	Code		string 	`gorm:"column:code" json:"code"`
	Category	string 	`gorm:"column:category" json:"category"`
	Icon		string 	`gorm:"column:icon" json:"icon"`
	Status		string 	`gorm:"column:status" json:"status"`
}

func GetFpxList() ([]*PaymentGatewaySetup, error) {
	var payment []*PaymentGatewaySetup
	db := GetDB()
	query := db.Table("payment_gateway_setup as a").
		Where("a.status = ?", "A")

	err := query.Order("id").Select("a.*").Find(&payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}