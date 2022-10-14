package model

import "gorm.io/gorm"

type CategorySetupStruct struct {
	Id           int64  `gorm:"column:id" json:"-"`
	CategoryName string `gorm:"column:category_name" json:"category_name"`
	CategoryCode string `gorm:"column:category_code" json:"category_code"`
}

func InitTable() *gorm.DB {
	return db.Table("category_setup")
}

func GetCategorySetupByCode(code string) (*CategorySetupStruct, error) {
	var categorySetup *CategorySetupStruct
	query := db.Table("category_setup as a").
		Where("category_code = ?", code)
	err := query.Select("a.*").First(&categorySetup).Error
	if err != nil {
		return nil, err
	}
	return categorySetup, nil
}

func GetAllCategorySetup(categoryId string, categoryCode string, categoryType string) (categories *CategorySetupStruct, err error) {
	query := InitTable()
	if categoryCode != "" {
		query = query.Where("category_code = ?", categoryCode)
	}
	if categoryId != "" {
		query = query.Where("id = ?", categoryId)
	}
	if categoryType != "" {
		query = query.Where("category_type = ?", categoryType)
	}
	err = query.Find(&categories).Error
	if err != nil {
		return
	}
	return
}
