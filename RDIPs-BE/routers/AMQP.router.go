package routers

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	"RDIPs-BE/handler"
	AMQP_handler "RDIPs-BE/handler/AMQP"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"os"
	"strconv"

	"github.com/rabbitmq/amqp091-go"
)

// var m sync.Mutex

func InitAmqpRoutes() {
	utils.Log(LogConstant.Info, "Initialize AMQP routes")
	defer utils.Log(LogConstant.Info, "Finish initialize AMQP routes")
	amqpPool := handler.GetRabbitPool()
	ch, _, err := amqpPool.Get()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	channel, ok := ch.(commonModel.BaseAmqpChannel)

	if !ok {
		utils.Log(LogConstant.Fatal, "Wrong format")
	}

	channel.Qos(10, 0, false)

	queue, err := channel.QueueDeclare("API", true, false, false, false, amqp091.Table{
		"x-message-ttl": 600000,
	})
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	priority, err := strconv.Atoi(os.Getenv("CONSUMER_PRIORITY"))
	if err != nil {
		utils.Log(LogConstant.Error, err)
		priority = 0 // default priority
	}

	deliveries, err := channel.Consume(
		queue.Name, // name
		"",         // consumerTag,
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		amqp091.Table{
			"x-priority": priority,
		})
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
	_, err = channel.QueueDeclare("SERVER", true, false, false, false, amqp091.Table{
		"x-message-ttl": 600000,
	})
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
}
