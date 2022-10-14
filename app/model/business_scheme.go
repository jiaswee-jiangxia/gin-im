package model

type BusinessSchemeStruct struct {
	Id         int    `gorm:"column:id" json:"-"`
	SchemeName string `gorm:"column:scheme_name" json:"scheme_name"`
	Status     string `gorm:"column:status" json:"status"`
	HomeShow   string `gorm:"column:home_show" json:"home_show"`
}

func GetBusinessScheme() ([]*BusinessSchemeStruct, error) {
	var Scheme []*BusinessSchemeStruct
	query := db.Table("business_scheme as a").
		Where("a.home_show = ?", "A").
		Limit(3).Order("a.sort")
	err := query.Find(&Scheme).Error
	if err != nil {
		return nil, err
	}
	return Scheme, nil
}
