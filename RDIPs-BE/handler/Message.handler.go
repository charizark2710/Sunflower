package handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	connection "RDIPs-BE/handler/Connection"
	commonModel "RDIPs-BE/model/common"

	"RDIPs-BE/utils"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

var rabbitPool *connection.Pool

type MessageHandler interface {
	Send(exchange string, body interface{}, deliveryMode uint8, replyTo string,
		correlationID string, routingKeyArgs ...string) error
}

type messageStruct struct{}

func NewMessageHandler() MessageHandler {
	return &messageStruct{}
}

func (m *messageStruct) Send(exchange string, body interface{}, deliveryMode uint8, replyTo string, correlationID string, routingKeyArgs ...string) error {
	routingKey := m.generateRoutingKey(routingKeyArgs...)
	utils.Log(LogConstant.Info, "Sending message to ", exchange, "with ", routingKey)
	conn, ctx, err := rabbitPool.Get()
	if err != nil {
		utils.Log(LogConstant.Error, err)
	} else {
		defer rabbitPool.Release(conn, ctx)
		channel, ok := conn.(commonModel.BaseAmqpChannel)
		if !ok {
			utils.Log(LogConstant.Error, "wrong channel format")
			return fmt.Errorf("wrong channel format")
		}
		var message []byte
		message, err = json.Marshal(body)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			return err
		}
		err = channel.PublishWithContext(context.Background(), exchange, routingKey, true, false, amqp091.Publishing{
			DeliveryMode: deliveryMode,
			ContentType:  "text/plain",
			Body:         message,
			ReplyTo:      correlationID,
			Headers:      amqp091.Table{"Correlation-ID": correlationID},
		})
		utils.Log(LogConstant.Info, "Finish sending message to ", exchange, "with ", routingKey)

	}
	return err
}

func (*messageStruct) generateRoutingKey(args ...string) string {
	args = append(args, "")
	copy(args[1:], args)
	args[0] = "server"
	return strings.Join(args, ".")
}

func GetRabbitPool() *connection.Pool {
	return rabbitPool
}

func SetRabbitPool(pool *connection.Pool) {
	rabbitPool = pool
}
