package model

import (
	"gorm.io/gorm"
)

type helper struct {
	db *gorm.DB
}

var Helper *helper = &helper{}

func (h *helper) GetDb() *gorm.DB {
	return h.db
}

func (h *helper) SetDb(db *gorm.DB) {
	h.db = db
}
