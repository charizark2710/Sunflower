package utils

import (
	"fmt"

	"github.com/charizark2710/Sunflower/RDIPs-BE/constant"
)

var logList = make(map[string]string)

func Log(id string, arg ...string) {
	fmt.Printf(logList[id], arg)
}

func PrepareLog() {
	logList = map[string]string{
		constant.Info:    "This is info log %s",
		constant.Warning: "This is warning log %s",
		constant.Error:   "This is error log %s",
		constant.Debug:   "This is debug log %s",
	}
}
