package utils

import (
	"log"

	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
)

var logList = make(map[string]string)

func Log(id string, arg ...interface{}) {
	if id == LogConstant.Error {
		log.Fatalln(logList[id], arg)
	} else {
		log.Println(logList[id], arg)
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
