package main

import (
	LogConstant "RDIPs-Task/constant/LogConst"
	AMQP_Handler "RDIPs-Task/handler/AMQP"
	Redis_Handler "RDIPs-Task/handler/Redis"
	"RDIPs-Task/utils"
)

func main() {
	utils.PrepareLog()
	AMQP_Handler.Connect()
	Redis_Handler.ConnectRedis()
	utils.Log(LogConstant.Info, "Start listening")
	select {}
}
