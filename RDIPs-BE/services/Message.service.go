package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

var TransmitMessage = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "TransmitMessage Start")
	defer utils.Log(LogConstant.Info, "TransmitMessage End")

	deviceID := c.Param("id")
	var payloadMSG interface{}
	err := json.Unmarshal(c.Body, &payloadMSG)
	if err == nil {
		err := handler.NewDeviceHandler(c.Ctx, &model.SysDevices{}).ReadDetail(false, deviceID)
		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
		}
		messageHandler := handler.NewMessageHandler()
		messageHandler.Send("amq."+amqp091.ExchangeTopic, payloadMSG, uuid.NewString(), deviceID)
		return commonModel.ResponseTemplate{HttpCode: 200, Data: "Transmit successfully"}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}
