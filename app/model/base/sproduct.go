package base

type SProductModel struct {
	Id			int64 `gorm:"column:id" json:"id"`
	Name		string `gorm:"column:name" json:"name"`
	Key			string `gorm:"column:key" json:"key"`
	Comment		string `gorm:"column:comment" json:"comment"`
	Updatetime	string `gorm:"column:updatetime" json:"-"`
	Createtime	string `gorm:"column:createtime" json:"-"`
	Status		string `gorm:"column:status" json:"status"`
	Sort		string `gorm:"column:sort" json:"sort"`
}

func GetProductByName(productName string) (*SProductModel,error) {

	return nil, nil
}