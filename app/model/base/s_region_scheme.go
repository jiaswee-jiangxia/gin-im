package base

type SRegionSchemeStruct struct {
	Id				int64 `gorm:"column:id" json:"-"`
	Name			string `gorm:"column:name" json:"name"`
	Type			int `gorm:"column:type" json:"type"`
	Updatetime		string `gorm:"column:updatetime" json:"updatetime"`
	Createtime		string `gorm:"column:createtime" json:"createtime"`
}

const LIVE_BLACK = 0

func FindSRegionSchemeById(schemeId int) (*SRegionSchemeStruct, error) {

	return nil, nil
}