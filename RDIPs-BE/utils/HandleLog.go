package utils

import (
	"log"

	LogConstant "RDIPs-BE/constant/LogConst"
)

var logList = make(map[string]string)

func Log(id string, arg ...interface{}) {
	switch id {
	case LogConstant.Fatal:
		log.Fatalf(logList[id], arg)
	default:
		log.Printf(logList[id], arg)
	}
}

func PrepareLog() {
	logList = map[string]string{
		LogConstant.Info:    "[INFO] %s",
		LogConstant.Warning: "[WARNING] %s",
		LogConstant.Error:   "[ERROR] %s",
		LogConstant.Debug:   "[DEBUG] %s",
	}
}
