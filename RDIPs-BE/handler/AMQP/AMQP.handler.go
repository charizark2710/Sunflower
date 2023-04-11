package AMQP_handler

import (
	AMQPconst "RDIPs-BE/constant/AMQP_Const"
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/constant/ServiceConst"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
			ack(d)
		} else {
			utils.Log(LogConstant.Info, "Finish ACK Delivery: "+d.Exchange+" With key: "+d.RoutingKey)
		}
	}
	defer cleanup()
	for delivery := range deliveries {
		utils.Log(LogConstant.Info, "Start Exchange: "+delivery.Exchange+" With key: "+delivery.RoutingKey)
		c := gin.Context{Request: &http.Request{
			Header: http.Header{"Content-Type": {"application/json"}, "Request-Type": {"amqp"}},
			Body:   http.NoBody,
		}}
		c.Request.Body = io.NopCloser(strings.NewReader(string(delivery.Body)))
		routingKeyArr := strings.Split(delivery.RoutingKey, ".")
		service := ServiceConst.ServicesMap[routingKeyArr[len(routingKeyArr)-1]]
		if fn, ok := service.(func(c *gin.Context) (commonModel.ResponseTemplate, error)); !ok {
			utils.Log(LogConstant.Error, "Wrong type of services functions")
		} else {
			result, err := fn(&c)
			if err != nil {
				utils.Log(LogConstant.Error, err)
				result.SetMessage(err.Error())
			}
			response, _ := json.Marshal(result)
			// if result.
			Send(AMQPconst.DATA_EXCHANGE, []string{"*"}, response)
		}
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
