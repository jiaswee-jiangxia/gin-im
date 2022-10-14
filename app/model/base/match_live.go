package base

import (
	_const "goskeleton/app/const"
	"goskeleton/app/model"
)

type MatchLiveStruct struct {
	Id				string `gorm:"column:id" json:"id"`
	MatchId			string `gorm:"column:match_id" json:"match_id"`
	MatchType		string `gorm:"column:match_type" json:"match_type"`
	Title			string `gorm:"column:title" json:"title"`
	Link			string `gorm:"column:link" json:"link"`
	Source			int `gorm:"column:source" json:"source"`
	Channel			string `gorm:"column:channel" json:"channel"`
	Name			string `gorm:"column:name" json:"name"`
	RegionSchemeId	int `gorm:"column:region_scheme_id" json:"region_scheme_id"`
}

func FindMatchLiveByProduct(matchType string, ids string, productId int64) ([]*MatchLiveStruct,error) {
	var mLiveData []*MatchLiveStruct
	db := model.GetDB()
	query1 := db.Table("match_live as a").
		Where("a.match_id = ?", ids).
		Where("a.match_type = ?", matchType).
		Where("a.source NOT IN (?,?)", _const.SOURCE_PARSE, _const.SOURCE_JUMP).
		Select("a.*")

		query2 := db.Table("match_live as a").
			Joins("inner join match_live_product b ON a.id = b.match_live_id").
			Where("b.product_id = ?", productId).
			Where("a.match_id = ?", ids).
			Where("a.match_type = ?", matchType).
			Select("a.*")

	err := db.Raw("? UNION ?",
		query1,
		query2,
	).Find(&mLiveData).Error

	if err != nil {
		return nil, err
	}
	return mLiveData, nil
}