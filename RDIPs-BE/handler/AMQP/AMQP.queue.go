package AMQP

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	connection "RDIPs-BE/handler/Connection"
	"RDIPs-BE/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

// var m sync.Mutex
var rabbitPool *connection.Pool = &connection.Pool{}

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

func (*messageStruct) generateRoutingKey(args ...string) string {
	return strings.Join(args, ".")
}

func GetRabbitPool() *connection.Pool {
	return rabbitPool
}

func SetRabbitPool(pool *connection.Pool) {
	rabbitPool = pool
}
