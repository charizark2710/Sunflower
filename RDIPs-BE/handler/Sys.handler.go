package handler

import (
	"reflect"

	commonModel "github.com/charizark2710/Sunflower/RDIPs-BE/model/common"
	"gorm.io/gorm"
)

var db = commonModel.DbHelper.GetDb()

func getStructType(i interface{}, db *gorm.DB) interface{} {
	reflectModel := reflect.ValueOf(i)
	m := reflectModel.Interface()
	if !db.Migrator().HasTable(m) {
		err := db.Migrator().CreateTable(m)
		if err != nil {
			return err
		}
		return err
	}
	return m
}

func Read(i interface{}) error {
	model := getStructType(i, db)
	err := db.Find(model).Where("deleted = ?", false).Error
	return err
}

func Create(i interface{}) error {
	m := getStructType(i, db)
	return db.Save(m).Error
}

func ReadDetail(i interface{}, id string) error {
	model := getStructType(i, db)
	err := db.Where("id = ? AND deleted = ?", id, false).First(model).Error
	return err
}

func Update(i interface{}, updatedData interface{}) error {
	model := getStructType(i, db)
	err := db.Model(model).Updates(updatedData).Error
	return err
}
