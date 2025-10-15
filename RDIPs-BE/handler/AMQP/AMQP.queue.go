package AMQP_handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	"RDIPs-BE/utils"
	"os"
	"strconv"

	"github.com/rabbitmq/amqp091-go"
)

// var m sync.Mutex

func InitAmqpQueue(channel *amqp091.Channel) error {
	utils.Log(LogConstant.Info, "Initialize AMQP routes")
	defer utils.Log(LogConstant.Info, "Finish initialize AMQP routes")
	channel.Qos(10, 0, false)

	queue, err := channel.QueueDeclare("API", true, false, false, false, amqp091.Table{
		"x-message-ttl": 600000,
	})
	if err != nil {
		return err
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
		return err
	}
	go func() {
		ReceiveService(deliveries)
	}()

	for serviceName := range ServiceConst.ServiceMapMQTT {
		err = channel.QueueBind(queue.Name, "gateway.*."+serviceName, "amq."+amqp091.ExchangeTopic, false, nil)
		if err != nil {
			return err
		}
		utils.Log(LogConstant.Debug, "Start Binding "+serviceName)
	}

	// declare queue for server
	_, err = channel.QueueDeclare("SERVER", true, false, false, false, amqp091.Table{
		"x-message-ttl": 600000,
	})
	if err != nil {
		return err
	}

	return nil
}
