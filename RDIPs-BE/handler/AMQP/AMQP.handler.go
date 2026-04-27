package AMQP

import (
	"RDIPs-BE/constant"
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	"RDIPs-BE/handler"
	connection "RDIPs-BE/handler/Connection"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"encoding/json"
	"errors"
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
	}

	utils.Log(LogConstant.Info, "Finish Connect Rabbitmq, err:", err)
	return err
}

func ReceiveService(deliveries <-chan amqp091.Delivery) {
	utils.Log(LogConstant.Info, "Start Receiver")
	var ack func(d *amqp091.Delivery, sysErr error)
	ack = func(d *amqp091.Delivery, sysErr error) {
		utils.Log(LogConstant.Info, "Start ACK Delivery: "+d.Exchange+" With key: "+d.RoutingKey)
		err := d.Ack(false)
		if err != nil {
			utils.Log(LogConstant.Error, err)
			time.Sleep(10 * time.Second)
			ack(d, sysErr)
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
		header[constant.REQUEST_TYPE_HEADER] = []string{"amqp"}
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
		var sysErr error
		if !ok {
			utils.Log(LogConstant.Error, "Service "+routingKeyArr[len(routingKeyArr)-1]+" is not exist")
			response = map[string]string{"ERROR": "Service" + routingKeyArr[len(routingKeyArr)-1] + "is not exist"}
		} else {
			result, err := fn(&c)
			if err != nil {
				sysErr = err
				utils.Log(LogConstant.Error, err)
				result.Error = err
				result.SetMessage(err.Error())
			}
			response = result
		}
		messageHandler := handler.NewMessageHandler()
		resBody, err := getResBody(response)
		if err == nil && resBody["needResponse"] == true {
			// After delete, optimize response
			if len(header["Correlation-Id"]) != 0 {
				deliveryMode, ok := resBody["deliveryMode"].(uint8)
				if !ok {
					deliveryMode = amqp091.Transient
				}
				// Response to the request device
				go messageHandler.Send(delivery.Exchange, response, deliveryMode, header["Correlation-Id"][0], "", routingKeyArr[1])
			}
		}

		go ack(&delivery, sysErr)
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
		utils.Log(LogConstant.Error, body, err)
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
		queries, ok := res["query"].(map[string]interface{})
		if ok {
			for key, value := range queries {
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

	if res["header"] != nil {
		headers, ok := res["header"].(map[string]interface{})
		if ok {
			for key, value := range headers {
				v := fmt.Sprintf("%v", value)
				c.Header.Add(key, v)
			}
		}
	}

}

func getResBody(response interface{}) (map[string]interface{}, error) {
	res, err := json.Marshal(response)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	resBody := make(map[string]interface{})
	err = json.Unmarshal(res, &resBody)
	if err != nil {
		utils.Log(LogConstant.Warning, "Cannot unmarshal delivery body: ", err)
		return nil, err
	}
	return resBody, err
}
