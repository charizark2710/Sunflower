package AMQP_handler

import (
	AMQPconst "RDIPs-BE/constant/AMQP_Const"
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func initializeAMQP() (commonModel.BaseAmqpConn, error) {
	conn, err := amqp091.Dial(
		"amqp://" +
			os.Getenv("BROKER_USER") +
			":" + os.Getenv("BROKER_PASSWORD") +
			"@" + os.Getenv("BROKER_HOST") +
			":" + os.Getenv("BROKER_PORT") + "/")
	if err != nil {
		return nil, err
	}

	return conn, err

	// // Declare exchange
	// for _, exchange := range AMQPconst.ExhangeArr {
	// 	ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	// }
}

func Send(exchange string, routingKeyArgs []string, body []byte) error {
	routingKey := generateRoutingKey(routingKeyArgs...)
	conn := commonModel.Helper.GetAMQPConnection()
	channel, err := (*conn).Channel()
	if err != nil {
		utils.Log(LogConstant.Error, err)
	} else {
		err = channel.PublishWithContext(context.Background(), exchange, routingKey, true, false, amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
		channel.Close()
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
		Send(AMQPconst.DATA_EXCHANGE, []string{"*"}, response)

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
