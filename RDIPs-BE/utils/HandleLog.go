package utils

import (
	"fmt"
	"log"
	"runtime"

	LogConstant "RDIPs-BE/constant/LogConst"
)

var logList = make(map[string]string)

func Log(id string, arg ...interface{}) {
	_, file, lineNo, _ := runtime.Caller(1)
	content := fmt.Sprintf("%v:%v::%v", file, lineNo, arg)
	switch id {
	case LogConstant.Fatal:
		log.Fatalf(logList[id], content)
	default:
		log.Printf(logList[id], content)
	}
}

func PrepareLog() {
	logList = map[string]string{
		LogConstant.Info:    "[INFO] %v",
		LogConstant.Warning: "[WARNING] %v",
		LogConstant.Error:   "[ERROR] %v",
		LogConstant.Debug:   "[DEBUG] %v",
	}
}
