package main

import (
	AMQP_Handler "RDIPs-Task/handler/AMQP"
	Redis_Handler "RDIPs-Task/handler/Redis"
)

func main() {
	AMQP_Handler.Connect()
	Redis_Handler.ConnectRedis()
	select {}
}
