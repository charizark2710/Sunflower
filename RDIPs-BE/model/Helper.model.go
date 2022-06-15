package model

import "gorm.io/gorm"

type helper struct {
	DB *gorm.DB
}

func (h *helper) GetDb() *gorm.DB {
	return h.DB
}

func (h *helper) SetDb(db *gorm.DB) {
	h.DB = db
}
