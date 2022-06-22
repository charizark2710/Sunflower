package model

type Devices struct {
	Id   string `json:"Id"`
	Name string `json:"name"`
}

type SysDevice struct {
	Id   string `gorm:"default:uuid_generate_v4();primaryKey;column:id;type:uuid"`
	Name string `gorm:"column:name"`
}

// func (SysDevice) TableName() string {
// 	return "sys_devices"
// }

// func (m *SysDevice) CreateDevices() error {
// 	db := DbHelper.GetDb()
// 	if !db.Migrator().HasTable(m) {
// 		err := db.Migrator().CreateTable(m)
// 		if err != nil {
// 			return err
// 		}
// 		err = db.Save(m).Error
// 		return err
// 	}
// 	return db.Save(m).Error
// }
