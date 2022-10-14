package base

type MRegionSchemeStruct struct {
	Id				int64 `gorm:"column:id" json:"-"`
	SchemeId		int `gorm:"column:scheme_id" json:"scheme_id"`
	RegionId		int `gorm:"column:region_id" json:"region_id"`
	RegionName		string `gorm:"column:region_name" json:"region_name"`
	Createtime		string `gorm:"column:createtime" json:"createtime"`
}

func FindMRegionSchemeById(schemeId int) ([]*MRegionSchemeStruct, error) {
	return nil, nil
}