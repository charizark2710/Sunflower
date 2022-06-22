package handler

import (
	"reflect"

	"github.com/charizark2710/Sunflower/RDIPs-BE/model"
)

func getStructType(i interface{}) interface{} {
	reflectModel := reflect.ValueOf(i)
	return reflectModel.Interface()
}

func Read(i interface{}) error {
	db := model.DbHelper.GetDb()
	model := getStructType(i)
	err := db.Find(model).Error
	return err
}

func Create(i interface{}) error {
	db := model.DbHelper.GetDb()
	m := getStructType(i)
	if !db.Migrator().HasTable(m) {
		err := db.Migrator().CreateTable(m)
		if err != nil {
			return err
		}
		err = db.Save(m).Error
		return err
	}
	return db.Save(m).Error
}

func ReadDetail(i interface{}, id string) error {
	db := model.DbHelper.GetDb()
	model := getStructType(i)
	err := db.Where("id = ?", id).First(model).Error
	return err
}
