package base

type BlackIpStruct struct {
	Id				int `gorm:"column:id" json:"-"`
	Ip				string `gorm:"column:ip" json:"ip"`
	Status			string `gorm:"column:status" json:"status"`
	Description		string `gorm:"column:description" json:"description"`
	IpIntStart		uint32 `gorm:"column:ip_int_start" json:"ip_int_start"`
	IpIntEnd		uint32 `gorm:"column:ip_int_end" json:"ip_int_end"`
}

const STATUS_OPEN = 1

func GetBlackIpList() ([]*BlackIpStruct, error) {
	return nil, nil
}
