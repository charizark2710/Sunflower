package AMQP_Handler

import (
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func Connect() (*amqp091.Connection, error) {
	amqpConn, err := amqp091.Dial(
		os.Getenv("BROKER_PROTOCOL") + "://" +
			os.Getenv("BROKER_USER") +
			":" + os.Getenv("BROKER_PASSWORD") +
			"@" + os.Getenv("BROKER_HOST") +
			":" + os.Getenv("BROKER_PORT") + "/")
	if err != nil {
		return nil, err
	}
	return amqpConn, nil
}

func ReceiveService(deliveries <-chan amqp091.Delivery) {

	var ack func(d *amqp091.Delivery, sysErr error)
	ack = func(d *amqp091.Delivery, sysErr error) {
		if err := d.Ack(false); err != nil {
			time.Sleep(10 * time.Second)
			ack(d, sysErr)
		}
	}

	// Map to route responses to correct goroutine
	for delivery := range deliveries {
		// routingKey := delivery.RoutingKey
		// send data to RL model
		// if score is larger than 0.3 then skip
		// else send to v8go to execute code
		// send back response
		go ack(&delivery, nil)
	}

}
