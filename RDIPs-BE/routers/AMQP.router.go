package routers

import (
	AMQPconst "RDIPs-BE/constant/AMQP_Const"
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	AMQP_handler "RDIPs-BE/handler/AMQP"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
)

// var m sync.Mutex

func InitAmqpRoutes() {
	utils.Log(LogConstant.Info, "Initialize AMQP routes")
	defer utils.Log(LogConstant.Info, "Finish initialize AMQP routes")
	ch := commonModel.Helper.GetAMQPChannel()

	queue, err := ch.QueueDeclare("API", true, false, false, false, nil)
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	deliveries, err := ch.Consume(
		queue.Name, // name
		"",         // consumerTag,
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
	go func() {
		AMQP_handler.ReceiveService(deliveries)
	}()
	for serviceName, _ := range ServiceConst.ServicesMap {
		err = ch.QueueBind(queue.Name, "gateway.*."+serviceName, AMQPconst.DATA_EXCHANGE, false, nil)
		if err != nil {
			utils.Log(LogConstant.Fatal, err)
		}
	}
}
