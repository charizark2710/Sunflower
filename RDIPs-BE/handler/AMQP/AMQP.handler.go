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
	"strings"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func Send(exchange string, routingKeyArgs []string, body []byte) error {
	routingKey := generateRoutingKey(routingKeyArgs...)
	channel, err := commonModel.Helper.GetAMQPConnection().Channel()
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
		c.Body = delivery.Body
		c.InitParamsAndQueries()
		c.SetQuery("amqp", "true")
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

func generateRoutingKey(args ...string) string {
	args = append(args, "")
	copy(args[1:], args)
	args[0] = "server"
	return strings.Join(args, ".")
}
