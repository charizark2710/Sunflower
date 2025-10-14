package AMQP_Handler

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"

	LogConstant "RDIPs-Task/constant/LogConst"
	handler "RDIPs-Task/handler"
	"RDIPs-Task/utils"
)

var receiveChannel *amqp091.Channel
var sendChannel *amqp091.Channel

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

	receiveChannel, err = amqpConn.Channel()
	sendChannel, err = amqpConn.Channel()
	if err != nil {
		return nil, err
	}
	return amqpConn, nil
}

func generateRoutingKey(args ...string) string {
	return strings.Join(args, ".")
}

func ReceiveService(deliveries <-chan amqp091.Delivery) {
	m := NewMessageHandler(sendChannel, receiveChannel)

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
			m.Send(delivery.ReplyTo, result, 1, delivery.CorrelationId, delivery.ReplyTo)
		} else {
			utils.Log(LogConstant.Error, err)
		}
		go ack(&delivery, nil)
	}

}
