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

	pingFn := func(conn interface{}) {
		ch, ok := conn.(commonModel.BaseAmqpChannel)
		if !ok {
			return
		}

		amqpClose := make(chan *amqp091.Error)
		ch.NotifyClose(amqpClose)

		go func() {
			utils.Log(LogConstant.Error, <-amqpClose)
		}()

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
	utils.Log(LogConstant.Info, "Sending message to ", exchange, "with ", routingKey)
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
		utils.Log(LogConstant.Info, "Finish sending message to ", exchange, "with ", routingKey)

	}
	return err
}

func ReceiveService(deliveries <-chan amqp091.Delivery) {
	utils.Log(LogConstant.Info, "Start Receiver")
	var ack func(d *amqp091.Delivery)
	ack = func(d *amqp091.Delivery) {
		utils.Log(LogConstant.Info, "Start ACK Delivery: "+d.Exchange+" With key: "+d.RoutingKey)
		err := d.Ack(false)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			time.Sleep(10 * time.Second)
			ack(d)
		} else {
			utils.Log(LogConstant.Info, "Finish ACK Delivery: "+d.Exchange+" With key: "+d.RoutingKey)
		}
	}
	defer func() {
		utils.Log(LogConstant.Info, "handle: deliveries channel closed")
	}()
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
		c.InitParamsAndQueries()
		c.SetQuery("amqp", "true")
		setGinContext(&c, delivery.Body)
		routingKeyArr := strings.Split(delivery.RoutingKey, ".")
		fn, ok := ServiceConst.ServicesMap[ServiceConst.ServiceMapMQTT[routingKeyArr[len(routingKeyArr)-1]]]
		var response []byte
		if !ok {
			utils.Log(LogConstant.Error, "Service", routingKeyArr[len(routingKeyArr)-1], "is not exist")
			response = []byte("Service" + routingKeyArr[len(routingKeyArr)-1] + "is not exist")
		} else {
			result, err := fn(&c)
			if err != nil {
				utils.Log(LogConstant.Error, err)
				result.Error = err
				result.SetMessage(err.Error())
			}
			response, err = json.Marshal(result)
			if err != nil {
				utils.Log(LogConstant.Error, err)
			}
		}
		// if result.
		go Send(delivery.Exchange, []string{string(c.CorrelationID())}, response)

		go ack(&delivery)
		utils.Log(LogConstant.Info, "Finish Exchange: "+delivery.Exchange+" With key: "+delivery.RoutingKey)
	}
	utils.Log(LogConstant.Info, "Done")
}

func setGinContext(c *commonModel.ServiceContext, body []byte) {
	res := make(map[string]interface{})
	err := json.Unmarshal(body, &res)
	if err != nil {
		utils.Log(LogConstant.Info, err)
		return
	}
	if res["param"] != nil {
		params, ok := res["param"].(map[string]string)
		if ok {
			for key, value := range params {
				c.SetParam(key, value)
			}
		}
	}

	if res["correlationID"] != nil {
		correlationID, ok := res["correlationID"].(string)
		if ok {
			c.SetCorrelationID(correlationID)
		}
	}

	if res["query"] != nil {
		querys, ok := res["query"].(map[string]string)
		if ok {
			for key, value := range querys {
				c.SetQuery(key, value)
			}
		}
	}

	if res["body"] != nil {
		body, marshalErr := json.Marshal(res["body"])
		if marshalErr != nil {
			utils.Log(LogConstant.Error, marshalErr)
		} else {
			c.Body = append(c.Body, body...)
		}
	}
}

func generateRoutingKey(args ...string) string {
	args = append(args, "")
	copy(args[1:], args)
	args[0] = "server"
	return strings.Join(args, ".")
}
