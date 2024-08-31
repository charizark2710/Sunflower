package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/handler"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

var TransmitMessage = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "TransmitMessage Start")
	defer utils.Log(LogConstant.Info, "TransmitMessage End")

	deviceID := c.Param("id")
	deliveryModeStr := c.Query("deliveryMode")
	deliveryMode, err := strconv.Atoi(deliveryModeStr)
	if err != nil {
		deliveryMode = int(amqp091.Transient)
	}
	if deliveryMode < 1 || deliveryMode > 2 {
		return commonModel.ResponseTemplate{HttpCode: 400, Data: nil}, fmt.Errorf("deliveryMode only allow 1 (Transient) or 2 (Persistent)")
	}
	var payloadMSG interface{}
	err = json.Unmarshal(c.Body, &payloadMSG)
	if err == nil {
		var device model.SysDevices
		err := handler.NewDeviceHandler(c.Ctx, &device).ReadDetail(false, deviceID)
		if err != nil {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil, Message: err.Error()}, err
		}
		messageHandler := handler.NewMessageHandler()
		messageHandler.Send("amq."+amqp091.ExchangeTopic, payloadMSG, uint8(deliveryMode), uuid.NewString(), device.Name)
		return commonModel.ResponseTemplate{HttpCode: 200, Data: "Transmit successfully"}, nil
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
}
