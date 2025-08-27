package handler

import (
	"RDIPs-Task/constant"
	LogConstant "RDIPs-Task/constant/LogConst"
	AMQP_Handler "RDIPs-Task/handler/AMQP"
	"RDIPs-Task/utils"
	"fmt"
	"os"
	"strconv"

	"context"
	"encoding/json"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler interface {
	Send(exchange string, body interface{}, deliveryMode uint8,
		correlationID string, routingKeyArgs ...string) error
}

type messageStruct struct {
	sendChannel    *amqp091.Channel
	receiveChannel *amqp091.Channel
}

func NewMessageHandler(sendChannel, receiveChannel *amqp091.Channel) MessageHandler {
	return &messageStruct{
		sendChannel:    sendChannel,
		receiveChannel: receiveChannel,
	}
}

func (m *messageStruct) Send(exchange string, body interface{}, deliveryMode uint8, correlationID string, routingKeyArgs ...string) error {
	routingKey := m.generateRoutingKey(routingKeyArgs...)
	var message []byte
	message, err := json.Marshal(body)
	if err != nil {
		return err
	}
	err = m.sendChannel.PublishWithContext(context.Background(), exchange, routingKey, true, false, amqp091.Publishing{
		DeliveryMode: deliveryMode,
		ContentType:  "text/plain",
		Body:         message,
		Headers:      amqp091.Table{"Correlation-ID": correlationID},
	})
	return err
}

func (m *messageStruct) InitAmqpQueue(conn *amqp091.Connection) {
	queue, err := m.receiveChannel.QueueDeclare("Task-Gateway", true, false, false, false, amqp091.Table{
		"x-message-ttl": 600000,
		"x-queue-type":  "stream",
	})
	if err != nil {
		fmt.Print(err)
	}

	priority, err := strconv.Atoi(os.Getenv("CONSUMER_PRIORITY"))
	if err != nil {
		utils.Log(LogConstant.Error, err)
		priority = 0 // default priority
	}

	deliveries, err := m.receiveChannel.Consume(
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
		AMQP_Handler.ReceiveService(deliveries)
	}()

	for _, routingKey := range constant.ROUTING_KEY_PREFIX {
		err = m.receiveChannel.QueueBind(queue.Name, routingKey+queue.Name, "amq."+amqp091.ExchangeFanout, false, nil)
	}
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
}

func (*messageStruct) generateRoutingKey(args ...string) string {
	return strings.Join(args, ".")
}
