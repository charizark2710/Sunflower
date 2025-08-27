package AMQP

import (
	"RDIPs-Scheduler/constant"
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	Bundle "RDIPs-Scheduler/handler"
	connection "RDIPs-Scheduler/handler/Connection"
	"RDIPs-Scheduler/utils"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

	notifyConnCloseCh := amqpConn.NotifyClose(make(chan *amqp091.Error, 1))

	// Reconnect if connection is close
	go func() {
		closedErr := <-notifyConnCloseCh
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
		if len(notifyConnCloseCh) > 0 {
			close(notifyConnCloseCh)
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
					bundleSize, metricMap, err := Bundle.AstParser(string(code))
					if err != nil {
						utils.Log(LogConstant.Error, err)
						break
					}

					// Send request
					requestCorrelationId := uuid.New().String()
					respCh := make(chan map[string]interface{}, 1)

					mu.Lock()
					pendingResponses[requestCorrelationId] = respCh
					mu.Unlock()

					MessageChan <- map[string]interface{}{
						"bundleSize":    bundleSize,
						"metricMap":     metricMap,
						"correlationId": requestCorrelationId,
					}

					// Wait with timeout
					ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
					select {
					case resp := <-respCh:
						bundleSize = resp["bundleSize"].(float64)
						metricMap = resp["metricMap"].(map[string]float64)

						messageHandler.Send(
							constant.JSCODE_ROUTING_KEY+"AGENTS",
							map[string]interface{}{
								"bundleSize": bundleSize,
								"metricMap":  metricMap,
							},
							1,
							requestCorrelationId,
						)

					case <-ctx.Done():
						utils.Log(LogConstant.Warning, "timeout waiting for correlationId "+requestCorrelationId)
						// Cleanup pending entry
						mu.Lock()
						delete(pendingResponses, requestCorrelationId)
						close(respCh)
						mu.Unlock()
					}
					cancel()
				}
				break // exit prefix loop after match
			}
		}
		go ack(&delivery, nil)
	}

	utils.Log(LogConstant.Debug, "Done")
}
