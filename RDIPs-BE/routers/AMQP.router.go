package routers

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	AMQP_handler "RDIPs-BE/handler/AMQP"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"

	"github.com/rabbitmq/amqp091-go"
)

// var m sync.Mutex

func InitAmqpRoutes() {
	utils.Log(LogConstant.Info, "Initialize AMQP routes")
	defer utils.Log(LogConstant.Info, "Finish initialize AMQP routes")
	amqpPool := AMQP_handler.GetPool()
	ch, err := amqpPool.Get()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	channel, ok := ch.(commonModel.BaseAmqpChannel)

	if !ok {
		utils.Log(LogConstant.Fatal, "Wrong format")
	}

	queue, err := channel.QueueDeclare("API", true, false, false, false, nil)
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	deliveries, err := channel.Consume(
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

	for serviceName := range ServiceConst.ServiceMapMQTT {
		err = channel.QueueBind(queue.Name, "gateway.*."+serviceName, "amq."+amqp091.ExchangeTopic, false, nil)
		if err != nil {
			utils.Log(LogConstant.Fatal, err)
		}
		utils.Log(LogConstant.Debug, "Start Binding "+serviceName)
	}

	// declare queue for server
	queueServer, err := channel.QueueDeclare("SERVER", true, false, false, false, nil)
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
	err = channel.QueueBind(queueServer.Name, "server.*", "amq."+amqp091.ExchangeTopic, false, nil)
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
}
