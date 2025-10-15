package AMQP_Handler

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"

	LogConstant "RDIPs-Task/constant/LogConst"
	handler "RDIPs-Task/handler"
	"RDIPs-Task/utils"
)

var receiveChannel *amqp091.Channel
var sendChannel *amqp091.Channel
var messageHandler MessageHandler

func Connect() (*amqp091.Connection, error) {
	amqpConn, err := amqp091.Dial(
		os.Getenv("BROKER_PROTOCOL") + "://" +
			os.Getenv("BROKER_USER") +
			":" + os.Getenv("BROKER_PASSWORD") +
			"@" + os.Getenv("BROKER_HOST") +
			":" + os.Getenv("BROKER_PORT") + "/")
	if err != nil {
		return nil, err
	}

	amqpConn.Config.Heartbeat = 10 * time.Second
	notifyConnCloseConn := amqpConn.NotifyClose(make(chan *amqp091.Error, 1))

	// Reconnect if connection is close
	go func() {
		closedErr := <-notifyConnCloseConn
		if closedErr != nil {
			utils.Log(LogConstant.Error, closedErr)
			_, err := Connect()
			// Open new channel if it get error
			for err != nil {
				utils.Log(LogConstant.Error, err)
				time.Sleep(10 * time.Second)
				_, err = Connect()
			}
		}
		if len(notifyConnCloseConn) > 0 {
			close(notifyConnCloseConn)
		}
	}()

	receiveChannel, err = amqpConn.Channel()
	if err != nil {
		return nil, err
	}
	sendChannel, err = amqpConn.Channel()
	if err != nil {
		return nil, err
	}
	messageHandler = NewMessageHandler(sendChannel, receiveChannel)
	return amqpConn, nil
}

func ReceiveService(deliveries <-chan amqp091.Delivery) {
	var ack func(d *amqp091.Delivery, sysErr error)
	ack = func(d *amqp091.Delivery, sysErr error) {
		if err := d.Ack(false); err != nil {
			time.Sleep(10 * time.Second)
			ack(d, sysErr)
		}
	}

	for delivery := range deliveries {
		var msg map[string]interface{}
		json.Unmarshal(delivery.Body, &msg)
		id := delivery.CorrelationId
		result, err := handler.ExecutionHandler(id, msg)
		if err == nil {
			messageHandler.Send(delivery.ReplyTo, result, 1, delivery.CorrelationId, delivery.ReplyTo)
		} else {
			utils.Log(LogConstant.Error, err)
		}
		go ack(&delivery, nil)
	}

}
