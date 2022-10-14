package model

type SysGeneralSetupStruct struct {
	Id    int    `gorm:"column:id" json:"-"`
	Name  string `gorm:"column:name" json:"name"`
	Value string `gorm:"column:value" json:"value"`
}

func GetSettingByName(name string) (*SysGeneralSetupStruct, error) {
	var Setting *SysGeneralSetupStruct
	db := GetDB()
	query := db.Table("sys_general_setup as a").
		Where("a.name = ?", name)
	err := query.First(&Setting).Error
	if err != nil {
		return nil, err
	}
	return Setting, nil
}
