package main

import (
	LogConstant "RDIPs-Task/constant/LogConst"
	AMQP_Handler "RDIPs-Task/handler/AMQP"
	Redis_Handler "RDIPs-Task/handler/Redis"
	"RDIPs-Task/utils"
)

func main() {
	utils.PrepareLog()
	_, err := AMQP_Handler.Connect()
	if err != nil {
		utils.Log(LogConstant.Fatal, err)
	}
	Redis_Handler.ConnectRedis()
	utils.Log(LogConstant.Info, "Start listening")
	select {}
}
