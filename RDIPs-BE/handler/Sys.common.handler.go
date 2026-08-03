package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	commonModel "RDIPs-BE/model/common"
)

type CommonHandler interface {
	Create() error
	Read(response interface{}, opts ...map[string]interface{}) error
	GetById(id string, response interface{}, opts ...map[string]interface{}) error
	Update() error
	Delete() error
}

type commonHandler struct {
	context    context.Context
	db         *gorm.DB
	mongoDB    *mongo.Database
	deviceBody interface{}
}

func newCommonHandler(c *gin.Context, dbType commonModel.DBType) CommonHandler {
	switch dbType {
	case commonModel.Postgres:
		return &commonHandler{db: GetDbFromContext[*gorm.DB](c, dbType), context: c}
	case commonModel.MongoDB:
		return &commonHandler{mongoDB: GetDbFromContext[*mongo.Database](c, dbType), context: c}
	default:
		return &commonHandler{db: GetDbFromContext[*gorm.DB](c, dbType), mongoDB: GetDbFromContext[*mongo.Database](c, dbType), context: c}
	}
}

func (*commonHandler) Read(interface{}, ...map[string]interface{}) error {
	return nil
}

func (*commonHandler) GetById(id string, response interface{}, opts ...map[string]interface{}) error {
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

func GetDbFromContext[T any](c *gin.Context, dbType commonModel.DBType) T {
	if c != nil {
		if val, exists := c.Get(string(dbType)); exists {
			if typedDB, ok := val.(T); ok {
				return typedDB
			}
		}
	}

	rawDB, err := commonModel.Factory.GetDB(dbType)
	if err != nil {
		panic("Failed to get database from factory: " + err.Error())
	}

	typedDB, ok := rawDB.(T)
	if !ok {
		panic("Failed to assert database type from factory")
	}

	return typedDB
}
