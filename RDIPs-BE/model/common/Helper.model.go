package model

import (
	"fmt"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var CacheSrv = memcache.New(os.Getenv("CACHE_SERVER"))

type DBType string

const (
	Postgres DBType = "postgres"
	MongoDB  DBType = "mongodb"
)

type DBFactory struct {
	databases map[DBType]any
}

var Factory = NewDBFactory()

func NewDBFactory() *DBFactory {
	return &DBFactory{
		databases: make(map[DBType]any),
	}
}

func (f *DBFactory) SetDB(dbType DBType, db any) {
	f.databases[dbType] = db
}

func (f *DBFactory) GetDB(dbType DBType) (any, error) {
	db, exists := f.databases[dbType]
	if !exists {
		return nil, fmt.Errorf("database driver %s not initialized", dbType)
	}
	return db, nil
}

func (f *DBFactory) GetGormDB() (*gorm.DB, error) {
	db, err := f.GetDB(Postgres)
	if err != nil {
		return nil, err
	}
	return db.(*gorm.DB), nil
}

func (f *DBFactory) GetMongoDB() (*mongo.Database, error) {
	db, err := f.GetDB(MongoDB)
	if err != nil {
		return nil, err
	}
	return db.(*mongo.Database), nil
}
