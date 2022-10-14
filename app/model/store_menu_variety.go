package model

type StoreMenuVarietyStruct struct {
	MenuId  string `gorm:"column:menu_id" json:"-"`
	Variety string `gorm:"column:variety" json:"variety"`
}

type VarietyListStruct struct {
	MaxChoices int                 `json:"max_choices"`
	MinChoices int                 `json:"min_choices"`
	Name       string              `json:"name"`
	Variety    []VarietyInfoStruct `json:"variety"`
}

type VarietyInfoStruct struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Required bool    `json:"required"`
}

func GetVarietyList(storeId int, menuId int) (*StoreMenuVarietyStruct, error) {
	var menu *StoreMenuVarietyStruct
	query := db.Table("store_menu_variety as a").
		Where("a.menu_id = ?", menuId)

	err := query.Select("a.*").First(&menu).Error
	if err != nil {
		return nil, err
	}
	return menu, nil
}
