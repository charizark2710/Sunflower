package AMQP_handler

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	"RDIPs-BE/handler"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var rabbitPool handler.Pool

func InitializeAMQP() error {
	conn, err := commonModel.Dial(
		"amqp://" +
			os.Getenv("BROKER_USER") +
			":" + os.Getenv("BROKER_PASSWORD") +
			"@" + os.Getenv("BROKER_HOST") +
			":" + os.Getenv("BROKER_PORT") + "/")
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
	factoryFn := func() (interface{}, error) {
		amqpCh, err := conn.Channel()
		return amqpCh, err
	}

	closeFn := func(conn interface{}) error {
		ch, ok := conn.(commonModel.BaseAmqpChannel)
		if !ok {
			return fmt.Errorf("%v", "wrong amqp connection format")
		}
		err := ch.Close()
		return err
	}

	pingFn := func(conn interface{}) chan interface{} {
		pingChan := make(chan interface{}, 1)
		ch, ok := conn.(commonModel.BaseAmqpChannel)
		if !ok {
			pingChan <- fmt.Errorf("%v", "wrong amqp connection format")
			return pingChan
		}

		if err != nil {
			pingChan <- err
			return pingChan
		}

		amqpErr := ch.NotifyClose(make(chan *amqp091.Error))

		go func() {
			pingChan <- amqpErr
		}()

		confirmCh := ch.NotifyPublish(make(chan amqp091.Confirmation))

		go func() {
			for {
				pingChan <- confirmCh
			}
		}()

		return pingChan
	}

	poolData := handler.PoolData{
		FactoryFn: factoryFn,
		CloseFn:   closeFn,
		PingFn:    pingFn,
	}

	err = rabbitPool.FillPool(poolData)

	return err
}

func GetPool() handler.Pool {
	return rabbitPool
}

func Send(exchange string, routingKeyArgs []string, body []byte) error {
	routingKey := generateRoutingKey(routingKeyArgs...)
	utils.Log(LogConstant.Info, "Sending message to %v with %v", exchange, routingKey)
	conn, err := rabbitPool.Get()
	if err != nil {
		utils.Log(LogConstant.Error, err)
	} else {
		defer rabbitPool.Release(conn)
		channel, ok := conn.(commonModel.BaseAmqpChannel)
		if !ok {
			return fmt.Errorf("wrong channel format")
		}
		err = channel.PublishWithContext(context.Background(), exchange, routingKey, true, false, amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
		utils.Log(LogConstant.Info, "Finish sending message to %v with %v", exchange, routingKey)

	}
	return err
}

func ReceiveService(deliveries <-chan amqp091.Delivery) {
	utils.Log(LogConstant.Info, "Start Receiver")
	cleanup := func() {
		utils.Log(LogConstant.Info, "handle: deliveries channel closed")
	}
	var ack func(d *amqp091.Delivery)
	ack = func(d *amqp091.Delivery) {
		utils.Log(LogConstant.Info, "Start ACK Delivery: "+d.Exchange+" With key: "+d.RoutingKey)
		err := d.Ack(true)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			time.Sleep(10 * time.Second)
			ack(d)
		} else {
			utils.Log(LogConstant.Info, "Finish ACK Delivery: "+d.Exchange+" With key: "+d.RoutingKey)
		}
	}
	defer cleanup()
	for delivery := range deliveries {
		utils.Log(LogConstant.Info, "Start Exchange: "+delivery.Exchange+" With key: "+delivery.RoutingKey)
		c := commonModel.ServiceContext{
			Ctx: context.Background(),
			Mu:  sync.Mutex{},
			ServiceModel: commonModel.ServiceModel{
				Body: []byte{},
				Header: http.Header{
					"Content-Type": {"application/json"}, "Request-Type": {"amqp"},
				},
			},
		}
		copy(c.Body, delivery.Body)
		c.InitParamsAndQueries()
		c.SetQuery("amqp", "true")
		setQueryAndParam(&c, delivery.Body)
		routingKeyArr := strings.Split(delivery.RoutingKey, ".")
		fn := ServiceConst.ServicesMap[routingKeyArr[len(routingKeyArr)-1]]
		result, err := fn(&c)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			result.Error = err
			result.SetMessage(err.Error())
		}
		response, _ := json.Marshal(result)
		// if result.
		Send(delivery.Exchange, []string{"*"}, response)

		go ack(&delivery)
		utils.Log(LogConstant.Info, "Finish Exchange: "+delivery.Exchange+" With key: "+delivery.RoutingKey)
	}
	utils.Log(LogConstant.Info, "Done")
}

func setQueryAndParam(c *commonModel.ServiceContext, body []byte) {
	res := make(map[string]interface{})
	err := json.Unmarshal(body, &res)
	utils.Log(LogConstant.Info, err)
	if err != nil && res["param"] != nil {
		params, ok := res["param"].(map[string]string)
		if ok {
			for key, value := range params {
				c.SetParam(key, value)
			}
		}
	}

	if err != nil && res["query"] != nil {
		querys, ok := res["query"].(map[string]string)
		if ok {
			for key, value := range querys {
				c.SetQuery(key, value)
			}
		}
	}
}

func generateRoutingKey(args ...string) string {
	args = append(args, "")
	copy(args[1:], args)
	args[0] = "server"
	return strings.Join(args, ".")
}
