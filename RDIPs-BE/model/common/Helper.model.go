package model

import (
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var CacheSrv = memcache.New(os.Getenv("CACHE_SERVER"))

type helper struct {
	db             *gorm.DB
	defaultMongoDB *mongo.Database
}

var Helper *helper = &helper{}

func (h *helper) GetDb() *gorm.DB {
	return h.db
}

func (h *helper) SetDb(db *gorm.DB) {
	h.db = db
}

func (h *helper) GetMongoDB() *mongo.Database {
	return h.defaultMongoDB
}

func (h *helper) SetMongoDB(mongoDB *mongo.Database) {
	h.defaultMongoDB = mongoDB
}
