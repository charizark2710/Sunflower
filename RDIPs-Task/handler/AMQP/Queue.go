package AMQP_Handler

import (
	"RDIPs-Task/constant"
	LogConstant "RDIPs-Task/constant/LogConst"
	"RDIPs-Task/utils"

	"context"
	"encoding/json"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler interface {
	Send(exchange string, body interface{}, deliveryMode uint8,
		correlationID string, routingKeyArgs ...string) error
	InitAmqpQueue() error
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

func (m *messageStruct) InitAmqpQueue() error {
	queue, err := m.receiveChannel.QueueDeclare(constant.EXECUTE_QUEUE, true, false, false, false, amqp091.Table{
		"x-max-age":    "5m",
		"x-queue-type": "stream",
	})
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}

	m.receiveChannel.Qos(10, 0, false)

	deliveries, err := m.receiveChannel.Consume(
		queue.Name, // name
		"",         // consumerTag,
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		amqp091.Table{})
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	go func() {
		ReceiveService(deliveries)
	}()

	for _, routingKey := range constant.ROUTING_KEY_POSTFIX {
		err = m.receiveChannel.QueueBind(queue.Name, queue.Name+"."+routingKey, "amq."+amqp091.ExchangeFanout, false, nil)
	}
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}

	return nil
}

func (*messageStruct) generateRoutingKey(args ...string) string {
	return strings.Join(args, ".")
}
