package utils

import (
	"fmt"

	LogConstant "github.com/charizark2710/Sunflower/RDIPs-BE/constant/LogConst"
)

var logList = make(map[string]string)

func Log(id string, arg ...interface{}) {
	fmt.Printf(logList[id], arg)
}

func PrepareLog() {
	logList = map[string]string{
		LogConstant.Info:    "This is info log %s",
		LogConstant.Warning: "This is warning log %s",
		LogConstant.Error:   "This is error log %s",
		LogConstant.Debug:   "This is debug log %s",
	}
}
