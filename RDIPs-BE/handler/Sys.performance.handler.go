package handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/model"
	"RDIPs-BE/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PerformanceHandler interface {
	CommonHandler
}

type performanceHandler struct {
	performanceBody *model.SysPerformance
	*commonHandler
}

func NewPerformanceHandler(c *gin.Context, performanceModel *model.SysPerformance) PerformanceHandler {
	commonHanlerInstance := newCommonHandler(c, "")
	commonStruct := commonHanlerInstance.(*commonHandler)
	return &performanceHandler{
		commonHandler:   commonStruct,
		performanceBody: performanceModel,
	}
}

func (p *performanceHandler) Read(performanceRes interface{}, opts ...map[string]interface{}) error {
	return p.db.Find(performanceRes).Error
}

func (p *performanceHandler) Create() error {
	performanceObj := p.performanceBody
	return p.db.Save(&performanceObj).Error
}

func (p *performanceHandler) GetById(id string, performanceResponse interface{}, opts ...map[string]interface{}) error {
	performance, ok := performanceResponse.(*model.SysPerformance)
	if !ok {
		return fmt.Errorf("wrong performance format")
	}
	err := p.db.Where("id = ?", id).First(performance).Error
	if err != nil {
		return err
	}

	if len(opts) > 0 {
		payload, err := p.handleOpts(performance.DocumentName, opts...)
		if err != nil {
			utils.Log(LogConstant.Error, err)
		} else {
			performance.Payload = payload
		}
	}
	return nil
}

func (p *performanceHandler) handleOpts(docName string, opts ...map[string]interface{}) (*[]map[string]interface{}, error) {
	// Generate pipeline for Aggregate
	mongoPipeLine := mongo.Pipeline{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		// limit stage
		limit, ok := opt["limit"]
		if ok && limit != "" {
			mongoPipeLine = append(mongoPipeLine, bson.D{
				{Key: "$limit", Value: limit},
			})
			continue
		}
		// filter stage
		filter, ok := opt["filterBytime"]
		if ok {
			filterMap, mapOk := filter.(map[string]time.Time)
			if !mapOk {
				return nil, fmt.Errorf("filter has wrong format")
			}
			from, fromOk := filterMap["from"]
			to, toOk := filterMap["to"]
			if !fromOk || !toOk {
				return nil, fmt.Errorf("from or to value is missing")
			}
			mongoPipeLine = append(mongoPipeLine, bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "timestamp", Value: bson.D{
						{Key: "$gte", Value: from}, {Key: "$lte", Value: to},
					}},
				}},
			})
			continue
		}
		for k, v := range opt {
			valueMap, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			value := p.handleCustomOpt(valueMap)
			if value == nil {
				continue
			}
			mongoPipeLine = append(mongoPipeLine, bson.D{
				{Key: "$match", Value: bson.D{
					{Key: k, Value: value},
				}},
			})
		}
	}
	// filter timestamp from mongodb
	payloadCursor, err := p.mongoDB.Collection(docName).Aggregate(context.TODO(), mongoPipeLine)

	if err != nil {
		return nil, err
	}
	return p.bsonToJson(payloadCursor)
}

func (p *performanceHandler) handleCustomOpt(valueMap map[string]interface{}) bson.D {
	var result bson.D
	for k, v := range valueMap {
		result = append(result, bson.E{
			Key: "$" + k, Value: v,
		})
	}
	return result
}

func (p *performanceHandler) bsonToJson(cursor *mongo.Cursor) (*[]map[string]interface{}, error) {
	var payload []map[string]interface{}
	for cursor.Next(context.Background()) {
		var bsonDoc bson.D
		err := cursor.Decode(&bsonDoc)
		if err != nil {
			return nil, err
		}
		temporaryBytes, err := bson.MarshalExtJSON(bsonDoc, true, true)
		if err != nil {
			return nil, err
		}
		var jsonDocument map[string]interface{}
		err = json.Unmarshal(temporaryBytes, &jsonDocument)
		if err != nil {
			return nil, err
		}
		payload = append(payload, jsonDocument)
	}
	return &payload, nil
}

func (p *performanceHandler) Update() error {
	var tempPerf model.SysPerformance
	err := p.GetById(p.performanceBody.Id, &tempPerf)
	if err != nil {
		utils.Log(LogConstant.Error, "Cannot find performance with Id = "+p.performanceBody.Id, err)
		return err
	}
	if tempPerf.DocumentName != p.performanceBody.DocumentName {
		err := fmt.Sprintf("Wrong documentName, documentName not exist: %s", p.performanceBody.DocumentName)
		utils.Log(LogConstant.Error, err)
		return fmt.Errorf(err)
	}
	mongoPayload := p.performanceBody.Payload
	if mongoPayload != nil {
		var documentPayload []interface{}
		for _, payload := range *mongoPayload {
			payload["timestamp"] = time.Now()
			documentPayload = append(documentPayload, payload)
		}
		_, err := p.mongoDB.Collection(p.performanceBody.DocumentName).
			InsertMany(context.TODO(), documentPayload, &options.InsertManyOptions{})
		if err != nil {
			return err
		}
	}
	err = p.db.Updates(p.performanceBody).Error
	// timeFinish := time.Now()
	// // Rollback if error
	// defer func() {
	// 	if err != nil {
	// 		p.mongoDb.DeleteMany(p.context, bson.D{{
	// 			Key: "timestamp", Value: bson.D{{Key: "$lt", Value: timeFinish}},
	// 		}})
	// 	}
	// }()

	return err
}
