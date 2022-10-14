package model

type CollectionSetup struct {
	Id			int 	`gorm:"column:id" json:"-"`
	Type		string 	`gorm:"column:type" json:"type"`
	Secret		string 	`gorm:"column:secret" json:"secret"`
	Title		string 	`gorm:"column:title" json:"title"`
	Status		string 	`gorm:"column:status" json:"status"`
}

func GetActiveCollection(collectionType string) (*CollectionSetup, error) {
	var payment *CollectionSetup
	db := GetDB()
	query := db.Table("collection_setup as a").
		Where("a.type = ?", collectionType).
		Where("a.status = ?", "A")

	err := query.Order("id").Select("a.*").First(&payment).Error
	if err != nil {
		return nil, err
	}
	return payment, nil
}