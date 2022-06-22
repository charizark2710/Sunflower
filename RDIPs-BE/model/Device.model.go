package model

type Devices struct {
	Id   string `json:"Id"`
	Name string `json:"name"`
}

type SysDevice struct {
	Id   string `gorm:"default:uuid_generate_v4();primaryKey;column:id;type:uuid"`
	Name string `gorm:"column:name"`
}
