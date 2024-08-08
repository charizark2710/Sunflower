package handler

import (
	commonModel "RDIPs-BE/model/common"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type CommonHandler interface {
	Create() error
	Read(response interface{}) error
	GetById(id string, response interface{}) error
	Update() error
	Delete() error
}

type commonHandler struct {
	context context.Context
	db      *gorm.DB
	mongoDB *mongo.Database
}

func newCommonHandler(c *gin.Context) CommonHandler {
	return &commonHandler{db: GetDbFromContext(c), mongoDB: GetMongoDBFromContext(c), context: c}
}

func (*commonHandler) Read(interface{}) error {
	return nil
}

func (*commonHandler) GetById(id string, response interface{}) error {
	return nil
}

func (*commonHandler) Create() error {
	return nil

}

func (*commonHandler) Delete() error {
	return nil

}

func (*commonHandler) Update() error {
	return nil
}

func GetDbFromContext(c *gin.Context) *gorm.DB {
	if c != nil {
		val, _ := c.Get("DB")
		if val != nil {
			res, ok := val.(*gorm.DB)
			if ok {
				return res
			}
		}
	}
	return commonModel.Helper.GetDb()
}

func GetMongoDBFromContext(c *gin.Context) *mongo.Database {
	if c != nil {
		val, _ := c.Get("MongoDB")
		if val != nil {
			res, ok := val.(*mongo.Database)
			if ok {
				return res
			}
		}
	}
	return commonModel.Helper.GetMongoDB()
}
