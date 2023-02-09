package model

import "gorm.io/gorm"

type helper struct {
	DB *gorm.DB
}

var DbHelper *helper = &helper{}

func (h *helper) GetDb() *gorm.DB {
	return h.DB
}

func (h *helper) SetDb(db *gorm.DB) {
	h.DB = db
}
