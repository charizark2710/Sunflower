package main

import (
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	"RDIPs-Scheduler/handler"
	"RDIPs-Scheduler/handler/AMQP"
	"RDIPs-Scheduler/utils"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

func main() {
	args := os.Args[1:]
	currentPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if len(args) > 0 && args[0] == "train" {
		// fd, err := openPerfEv()
		// if err != nil {
		// 	panic(err)
		// }
		var dir_path string
		if len(args) > 1 {
			dir_path = args[1]
		} else {
			dir_path = "./javascript-algorithms-and-data-structures"
		}
		socketPath, _ := filepath.Abs(currentPath + "/training.sock")
		files := handler.LoadJSDir(dir_path, 100)
		if len(files) > 0 {
			conn := handler.Connect(socketPath)
			for _, file := range files {
				// Handle bundled code
				prev := debug.SetGCPercent(-1)
				_, count, cycle, err := handler.ExecuteJs(file["bundledCode"], 3)
				debug.SetGCPercent(prev)
				runtime.GC()
				if err == nil {
					bundleSize, metricsMap, err := handler.AstParser(file["bundledCode"])

					if err != nil {
						utils.Log(LogConstant.Error, err)
						continue
					}
					obj := map[string]interface{}{
						"name":       strings.Split(file["name"], "/")[len(strings.Split(file["name"], "/"))-1],
						"code":       file["bundledCode"],
						"ic":         count,
						"cycle":      cycle,
						"bundleSize": bundleSize,
					}
					for k, v := range metricsMap {
						obj[k] = v
					}

					b, _ := json.Marshal(obj)
					conn.Write(append(b, []byte("\nEND\n")...))
				} else {
					continue
				}

				// // Handle unbundled code
				// _, timeExec, _, err = handler.ExecuteJs(file["orgCode"])

				// if err == nil {
				// 	bundleSize, metricsMap, err := handler.AstParser(file["orgCode"])

				// 	if err != nil {
				// 		utils.Log(LogConstant.Error, err)
				// 		continue
				// 	}
				// 	obj := map[string]interface{}{
				// 		"name":       strings.Split(file["name"], "/")[len(strings.Split(file["name"], "/"))-1],
				// 		"code":       file["orgCode"],
				// 		"time_us":    timeExec,
				// 		"bundleSize": bundleSize,
				// 	}
				// 	for k, v := range metricsMap {
				// 		obj[k] = v
				// 	}

				// 	b, _ := json.Marshal(obj)
				// 	conn.Write(append(b, []byte("\nEND\n")...))
				// }
			}
		}
	} else {
		initErr := AMQP.InitializeAMQP()
		if initErr != nil {
			panic(initErr)
		}
		utils.Log(LogConstant.Info, "Start listening")
		select {}
	}

	// bundleSize, functionComplexity, conditionalComplexity, err := handler.AstParser(string(data))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Function Complexity: %f \n", functionComplexity)
	// fmt.Printf("Conditional Complexity: %f \n", conditionalComplexity)
	// fmt.Printf("Bundle Size: %f bytes\n", bundleSize)
}
