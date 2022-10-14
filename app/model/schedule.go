package model

type MatchStruct struct {
	Id				string `gorm:"column:id" json:"-"`
	MatchType		string `gorm:"column:match_type" json:"match_type"`
}

func GetFootballMatchByTimeRange(dateFrom string, dateTo string) ([]*MatchStruct, error) {
	var Matches []*MatchStruct
	db := GetDB()
	err := db.Table("football_match").
		Where("time >= ?", dateFrom).
		Where("time <= ?", dateTo).
		Where("state >= ?", 0).
		Order("time").
		Select("id, 'zuqiu' as match_type").
		Find(&Matches).Error
	if err != nil {
		return nil, err
	}
	return Matches, nil
}

func GetBasketballMatchByTimeRange(dateFrom string, dateTo string) ([]*MatchStruct, error) {
	var Matches []*MatchStruct
	db := GetDB()
	err := db.Table("basketball_match").
		Where("time >= ?", dateFrom).
		Where("time <= ?", dateTo).
		Where("state >= ?", 0).
		Order("time").
		Select("id, 'lanqiu' as match_type").
		Find(&Matches).Error
	if err != nil {
		return nil, err
	}
	return Matches, nil
}

func GetZongheMatchByTimeRange(dateFrom string, dateTo string) ([]*MatchStruct, error) {
	var Matches []*MatchStruct
	db := GetDB()
	err := db.Table("zonghe_match").
		Where("time >= ?", dateFrom).
		Where("time <= ?", dateTo).
		Where("state >= ?", 0).
		Order("time").
		Select("id, 'zonghe' as match_type").
		Find(&Matches).Error
	if err != nil {
		return nil, err
	}
	return Matches, nil
}