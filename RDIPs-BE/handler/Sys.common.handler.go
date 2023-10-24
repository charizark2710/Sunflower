package handler

import (
	commonModel "RDIPs-BE/model/common"

	"github.com/gin-gonic/gin"
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
	db *gorm.DB
}

func newCommonHandler(c *gin.Context) CommonHandler {
	return &commonHandler{db: GetDbFromContext(c)}
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

// func getStructType(i interface{}, db *gorm.DB) error {
// 	reflectModel := reflect.ValueOf(i)
// 	m := reflectModel.Interface()
// 	if !db.Migrator().HasTable(m) {
// 		return errors.New("table is not exist")
// 	}
// 	return nil
// }
