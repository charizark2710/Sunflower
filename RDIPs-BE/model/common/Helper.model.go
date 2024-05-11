package model

import (
	"github.com/bradfitz/gomemcache/memcache"
	"gorm.io/gorm"
)

var CacheSrv = memcache.New("10.0.0.1:11211")

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
