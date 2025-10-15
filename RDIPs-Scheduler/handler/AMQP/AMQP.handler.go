package AMQP

import (
	"RDIPs-Scheduler/constant"
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	handler "RDIPs-Scheduler/handler"
	connection "RDIPs-Scheduler/handler/Connection"
	"RDIPs-Scheduler/utils"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var MessageChan = make(chan map[string]interface{})
var ResultChan = make(chan map[string]interface{})

var rabbitPool *connection.Pool = &connection.Pool{}

func InitializeAMQP() error {
	amqpConn, err := amqp091.Dial(
		os.Getenv("BROKER_PROTOCOL") + "://" +
			os.Getenv("BROKER_USER") +
			":" + os.Getenv("BROKER_PASSWORD") +
			"@" + os.Getenv("BROKER_HOST") +
			":" + os.Getenv("BROKER_PORT") + "/")
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return err
	}
	amqpConn.Config.Heartbeat = 10 * time.Second
	notifyConnCloseConn := amqpConn.NotifyClose(make(chan *amqp091.Error, 1))

	// Reconnect if connection is close
	go func() {
		closedErr := <-notifyConnCloseConn
		if closedErr != nil {
			utils.Log(LogConstant.Error, closedErr)
			rabbitPool.Close()

			err := InitializeAMQP()
			// Open new channel if it get error
			for err != nil {
				utils.Log(LogConstant.Error, err)
				time.Sleep(10 * time.Second)
				err = InitializeAMQP()
			}
		}
		if len(notifyConnCloseConn) > 0 {
			close(notifyConnCloseConn)
		}
	}()

	factoryFn := func() (interface{}, error) {
		amqpCh, err := amqpConn.Channel()
		return amqpCh, err
	}

	closeFn := func(conn interface{}) error {
		ch, ok := conn.(*amqp091.Channel)
		if !ok {
			return fmt.Errorf("%v", "wrong amqp connection format")
		}
		err := ch.Close()
		return err
	}

	pingFn := func(conn interface{}) error {
		ch, ok := conn.(*amqp091.Channel)
		if !ok {
			return errors.New("wrong connection")
		}
		err := InitAmqpQueue(ch)
		if err != nil {
			return err
		}
		chClose := ch.NotifyClose(make(chan *amqp091.Error, 1))
		// Re-initialize channel if this one is closed due to some error
		go func() {
			closedErr := <-chClose
			if closedErr != nil {
				utils.Log(LogConstant.Error, closedErr)
				if amqpConn.IsClosed() {
					return
				}
				amqpCh, err := amqpConn.Channel()
				if err == nil {
					err = InitAmqpQueue(amqpCh)
				}
				// Open new channel if it get error
				for err != nil {
					utils.Log(LogConstant.Error, err)
					time.Sleep(10 * time.Second)
					amqpCh, err = amqpConn.Channel()
				}
				conn = amqpCh
			}
			if len(chClose) > 0 {
				close(chClose)
			}
		}()
		return nil
	}

	poolData := connection.PoolData{
		FactoryFn: factoryFn,
		CloseFn:   closeFn,
		PingFn:    pingFn,
	}

	err = rabbitPool.FillPool(poolData)

	if err != nil {
		amqpConn.Close()
		return err
	}

	SetRabbitPool(rabbitPool)
	utils.Log(LogConstant.Info, "Finish Connect Rabbitmq, err:", err)
	return err
}

func ReceiveService(deliveries <-chan amqp091.Delivery) {
	messageHandler := NewMessageHandler()
	utils.Log(LogConstant.Info, "Start Receiver")
	var ack func(d *amqp091.Delivery, sysErr error)
	ack = func(d *amqp091.Delivery, sysErr error) {
		utils.Log(LogConstant.Info, "Start ACK Delivery: "+d.Exchange+" With key: "+d.RoutingKey)
		if err := d.Ack(false); err != nil {
			if err == amqp091.ErrClosed {
				return
			}
			utils.Log(LogConstant.Error, err)
			time.Sleep(10 * time.Second)
			ack(d, sysErr)
		} else {
			utils.Log(LogConstant.Info, "Finish ACK Delivery: "+d.Exchange+" With key: "+d.RoutingKey)
		}
	}

	// Map to route responses to correct goroutine
	pendingResponses := make(map[string]chan map[string]interface{})
	var mu sync.Mutex

	// Response router goroutine
	go func() {
		for resp := range ResultChan {
			corrID, ok := resp["correlationId"].(string)
			if !ok {
				continue
			}
			mu.Lock()
			if ch, found := pendingResponses[corrID]; found {
				ch <- resp
				close(ch)
				delete(pendingResponses, corrID)
			}
			mu.Unlock()
		}
	}()

	defer utils.Log(LogConstant.Info, "handle: deliveries channel closed")

	for delivery := range deliveries {
		routingKey := delivery.RoutingKey
		utils.Log(LogConstant.Debug, "Received a message: Delivery: "+delivery.Exchange+" With key: "+routingKey)

		for _, prefixRoutingKey := range constant.ROUTING_KEY_PREFIX {
			if strings.HasPrefix(routingKey, prefixRoutingKey) {
				switch prefixRoutingKey {
				case constant.JSCODE_ROUTING_KEY:
					code := delivery.Body
					result, err := handler.GuessHandler(string(code))
					if err == nil && result != nil && result["error"] == nil {
						messageHandler.Send(constant.EXECUTE_QUEUE, map[string]any{
							"ic":             result["ic"],
							"cycle":          result["cycle"],
							"bundleSize":     result["bundleSize"],
							"ast_metric":     result["ast_metric"],
							"input_ids":      result["input_ids"],
							"attention_mask": result["attention_mask"],
							"code":           result["code"]}, amqp091.Persistent, delivery.CorrelationId, delivery.ReplyTo)
						go ack(&delivery, nil)
					} else {
						delivery.Nack(false, true)
					}
				default:
					go ack(&delivery, nil)
				}
			}
		}

		utils.Log(LogConstant.Debug, "Done")
	}
}
