package handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/model"
	"RDIPs-BE/utils"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type PerformanceHandler interface {
	CommonHandler
}

type performanceHandler struct {
	performanceCollection *mongo.Collection
	performanceBody       *model.SysPerformance
	*commonHandler
}

func NewPerformanceHandler(c *gin.Context, performanceModel *model.SysPerformance) PerformanceHandler {
	commonHanlerInstance := newCommonHandler(c)
	commonStruct := commonHanlerInstance.(*commonHandler)
	return &performanceHandler{
		commonHandler:         commonStruct,
		performanceBody:       performanceModel,
		performanceCollection: commonStruct.mongoDB.Collection("performance"),
	}
}

func (p *performanceHandler) Read(performanceRes interface{}) error {
	return p.db.Find(performanceRes).Error
}

func (p *performanceHandler) Create() error {
	performanceObj := p.performanceBody
	return p.db.Save(&performanceObj).Error
}

func (p *performanceHandler) GetById(id string, performanceResponse interface{}) error {
	return p.db.Where("id = ?", id).First(performanceResponse).Error
}

func (p *performanceHandler) Update() error {
	err := p.GetById(p.performanceBody.Id, &model.SysPerformance{})
	if err != nil {
		utils.Log(LogConstant.Error, "Cannot find performance with Id = "+p.performanceBody.Id, err)
		return err
	}
	mongoPayload := p.performanceBody.Payload
	if mongoPayload != nil {
		var documentPayload []interface{}
		for _, payload := range *mongoPayload {
			payload["timestamp"] = time.Now()
			payload["document_name"] = p.performanceBody.DocumentName
			documentPayload = append(documentPayload, payload)
		}
		_, err := p.performanceCollection.InsertMany(context.TODO(), documentPayload)
		if err != nil {
			return err
		}
	}
	err = p.db.Updates(p.performanceBody).Error
	// timeFinish := time.Now()
	// // Rollback if error
	// defer func() {
	// 	if err != nil {
	// 		p.performanceCollection.DeleteMany(p.context, bson.D{{
	// 			Key: "timestamp", Value: bson.D{{Key: "$lt", Value: timeFinish}},
	// 		}})
	// 	}
	// }()

	return err
}
