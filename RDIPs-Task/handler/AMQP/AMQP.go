package AMQP_Handler

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"

	LogConstant "RDIPs-Task/constant/LogConst"
	handler "RDIPs-Task/handler"
	"RDIPs-Task/utils"
)

var channel *amqp091.Channel

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

	channel, err = amqpConn.Channel()
	if err != nil {
		return nil, err
	}
	return amqpConn, nil
}

func Send(exchange string, body interface{}, deliveryMode uint8, correlationID string, replyTo string, routingKeyArgs ...string) error {
	routingKey := generateRoutingKey(routingKeyArgs...)
	utils.Log(LogConstant.Info, "Sending message to ", exchange, "with ", routingKey)
	var message []byte
	message, err := json.Marshal(body)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	err = channel.PublishWithContext(context.Background(), exchange, routingKey, true, false, amqp091.Publishing{
		DeliveryMode: deliveryMode,
		ContentType:  "text/plain",
		Body:         message,
		Headers:      amqp091.Table{"Correlation-ID": correlationID},
		ReplyTo:      replyTo,
	})
	utils.Log(LogConstant.Info, "Finish sending message to ", exchange, "with ", routingKey)
	return err
}

func generateRoutingKey(args ...string) string {
	return strings.Join(args, ".")
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
		code := msg["code"].(string)
		ic := msg["ic"].(string)
		id := msg["id"].(string)
		result, err := handler.ExecutionHandler(id, code, ic)
		if err == nil {
			Send(delivery.ReplyTo, result, 1, delivery.CorrelationId, delivery.ReplyTo)
		} else {
			utils.Log(LogConstant.Error, err)
		}
		go ack(&delivery, nil)
	}

}
