package AMQP

import (
	"RDIPs-Scheduler/constant"
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	connection "RDIPs-Scheduler/handler/Connection"
	"os"
	"strconv"

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
		defer rabbitPool.Release(conn, ctx)
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

func (m *messageStruct) InitAmqpQueue() {
	utils.Log(LogConstant.Info, "Initialize AMQP routes")
	defer utils.Log(LogConstant.Info, "Finish initialize AMQP routes")
	amqpPool := GetRabbitPool()
	ch, _, err := amqpPool.Get()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}

	channel, ok := ch.(*amqp091.Channel)

	if !ok {
		utils.Log(LogConstant.Fatal, "Wrong format")
	}

	channel.Qos(10, 0, false)

	queue, err := channel.QueueDeclare(constant.SCHEDULER_QUEUE, true, false, false, false, amqp091.Table{
		"x-message-ttl": 600000,
		"x-queue-type":  "stream",
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
		ReceiveService(deliveries)
	}()

	for _, routingKey := range constant.ROUTING_KEY_PREFIX {
		err = channel.QueueBind(queue.Name, routingKey+queue.Name, "amq."+amqp091.ExchangeFanout, false, nil)
	}
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
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
