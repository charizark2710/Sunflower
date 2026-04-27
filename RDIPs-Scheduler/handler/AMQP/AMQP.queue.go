package AMQP

import (
	"RDIPs-Scheduler/constant"
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	connection "RDIPs-Scheduler/handler/Connection"

	"RDIPs-Scheduler/utils"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

type MessageHandler interface {
	Send(exchange string, body interface{}, deliveryMode uint8,
		correlationID string, replyTo string, routingKeyArgs ...string) error
}

type messageStruct struct{}

func NewMessageHandler() MessageHandler {
	return &messageStruct{}
}

func (m *messageStruct) Send(exchange string, body interface{}, deliveryMode uint8, correlationID string, replyTo string, routingKeyArgs ...string) error {
	routingKey := m.generateRoutingKey(routingKeyArgs...)
	utils.Log(LogConstant.Info, "Sending message to ", exchange, "with ", routingKey)
	conn, ctx, err := rabbitPool.Get()
	if err != nil {
		utils.Log(LogConstant.Error, err)
	} else {
		defer rabbitPool.Release(conn, false, ctx)
		channel, ok := conn.(*amqp091.Channel)
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
			Headers:      amqp091.Table{"Correlation-ID": correlationID},
			ReplyTo:      replyTo,
		})
		utils.Log(LogConstant.Info, "Finish sending message to ", exchange, "with ", routingKey)

	}
	return err
}

func InitAmqpQueue(channel *amqp091.Channel) error {
	utils.Log(LogConstant.Info, "Initialize AMQP routes")
	defer utils.Log(LogConstant.Info, "Finish initialize AMQP routes")

	channel.Qos(10, 0, false)

	queue, err := channel.QueueDeclare(constant.SCHEDULER_QUEUE, true, false, false, false, amqp091.Table{
		"x-queue-type": "stream",
		"x-max-age":    "5m",
	})
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	err = channel.ExchangeDeclare(constant.SCHEDULER_QUEUE, amqp091.ExchangeFanout, true, false, false, false, nil)
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
		amqp091.Table{})
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
	go func() {
		ReceiveService(deliveries)
	}()

	for _, routingKey := range constant.ROUTING_KEY_POSTFIX {
		err = channel.QueueBind(queue.Name, queue.Name+"."+routingKey, constant.SCHEDULER_QUEUE, false, nil)
	}
	if err != nil {
		return err
	}
	return nil
}

func (*messageStruct) generateRoutingKey(args ...string) string {
	return strings.Join(args, ".")
}

func GetRabbitPool() *connection.Pool {
	return rabbitPool
}

func SetRabbitPool(pool *connection.Pool) {
	rabbitPool = pool
}
