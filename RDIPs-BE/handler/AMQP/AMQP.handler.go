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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

func Send(exchange string, body interface{}, correlationID string, routingKeyArgs ...string) error {
	routingKey := generateRoutingKey(routingKeyArgs...)
	utils.Log(LogConstant.Info, "Sending message to ", exchange, "with ", routingKey)
	conn, err := rabbitPool.Get()
	if err != nil {
		utils.Log(LogConstant.Error, err)
	} else {
		defer rabbitPool.Release(conn)
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
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         message,
			Headers:      amqp091.Table{"Correlation-ID": correlationID},
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
		header := tableToHttpHeader(delivery.Headers)
		header["Content-Type"] = []string{"application/json"}
		header["Request-Type"] = []string{"amqp"}
		c := commonModel.ServiceContext{
			Ctx: &gin.Context{},
			Mu:  sync.Mutex{},
			ServiceModel: commonModel.ServiceModel{
				Body:   []byte{},
				Header: header,
			},
		}
		c.InitParamsAndQueries()
		c.SetQuery("amqp", "true")
		setGinContext(&c, delivery.Body)
		routingKeyArr := strings.Split(delivery.RoutingKey, ".")
		fn, ok := ServiceConst.ServicesMap[ServiceConst.ServiceMapMQTT[routingKeyArr[len(routingKeyArr)-1]]]
		var response interface{}
		if !ok {
			utils.Log(LogConstant.Error, "Service"+routingKeyArr[len(routingKeyArr)-1]+"is not exist")
			response = map[string]string{"ERROR": "Service" + routingKeyArr[len(routingKeyArr)-1] + "is not exist"}
		} else {
			result, err := fn(&c)
			if err != nil {
				utils.Log(LogConstant.Error, err)
				result.Error = err
				result.SetMessage(err.Error())
			}
			response = result
		}

		// TODO: Delete else after gateway is implemented
		// After delete, optimize response
		if header["Correlation-Id"] != nil && len(header["Correlation-Id"]) != 0 {
			go Send(delivery.Exchange, response, header["Correlation-Id"][0], "*")
		} else {
			res, err := json.Marshal(response)
			if err != nil {
				utils.Log(LogConstant.Error, err)
			}
			body := make(map[string]interface{})
			err = json.Unmarshal(delivery.Body, &body)
			if err != nil {
				utils.Log(LogConstant.Warning, "Cannot unmarshal delivery body: ", err)
			}
			resBody := make(map[string]interface{})
			err = json.Unmarshal(res, &resBody)
			if err != nil {
				utils.Log(LogConstant.Warning, "Cannot unmarshal response body: ", err)
			} else {
				resBody["CorrelationId"] = body["CorrelationId"]
			}
			id, _ := body["CorrelationId"].(string)
			go Send(delivery.Exchange, resBody, id, "*")
		}

		go ack(&delivery)
		utils.Log(LogConstant.Info, "Finish Exchange: "+delivery.Exchange+" With key: "+delivery.RoutingKey)
	}
	utils.Log(LogConstant.Info, "Done")
}

func tableToHttpHeader(table amqp091.Table) http.Header {
	header := http.Header{}
	for k, v := range table {
		switch value := v.(type) {
		case string:
			header[k] = []string{value}
		case []string:
			header[k] = value
		case bool:
			header[k] = []string{strconv.FormatBool(value)}
		case int64:
			header[k] = []string{strconv.FormatInt(value, 10)}
		case float64:
			header[k] = []string{fmt.Sprintf("%f", value)}
		}
	}
	return header
}

func setGinContext(c *commonModel.ServiceContext, body []byte) {
	res := make(map[string]interface{})
	err := json.Unmarshal(body, &res)
	if err != nil {
		utils.Log(LogConstant.Info, err)
		return
	}
	if res["param"] != nil {
		params, ok := res["param"].(map[string]interface{})
		if ok {
			for key, value := range params {
				v := fmt.Sprintf("%v", value)
				c.SetParam(key, v)
			}
		}
	}

	if res["query"] != nil {
		querys, ok := res["query"].(map[string]interface{})
		if ok {
			for key, value := range querys {
				v := fmt.Sprintf("%v", value)
				c.SetQuery(key, v)
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
