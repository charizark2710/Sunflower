package handler

import (
	"errors"
	"reflect"

	commonModel "RDIPs-BE/model/common"

	"gorm.io/gorm"
)

func getStructType(i interface{}, db *gorm.DB) error {
	reflectModel := reflect.ValueOf(i)
	m := reflectModel.Interface()
	if !db.Migrator().HasTable(m) {
		return errors.New("table is not exist")
	}
	return nil
}

// func Read(i interface{}) error {
// 	model := getStructType(i, db)
// 	err := db.Find(model).Where("deleted = ?", false).Error
// 	return err
// }

func Create(i interface{}) error {
	var db = commonModel.Helper.GetDb()
	err := getStructType(i, db)
	if err != nil {
		return err
	}
	return db.Save(i).Error
}

// func ReadDetail(i interface{}, id string) error {
// 	model := getStructType(i, db)
// 	err := db.Where("id = ? AND deleted = ?", id, false).First(model).Error
// 	return err
// }

func Update(i interface{}, updatedData interface{}) error {
	var db = commonModel.Helper.GetDb()
	err := getStructType(i, db)
	if err != nil {
		return err
	}
	return db.Model(i).Updates(updatedData).Error
}

func CreateWithTx(i interface{}, tx *gorm.DB) error {
	var db = commonModel.Helper.GetDb()
	if err := getStructType(i, db); err != nil {
		return err
	}
	return tx.Create(i).Error
}
