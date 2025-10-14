package handler

import (
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	"RDIPs-Scheduler/utils"
	"encoding/json"
	"io"
	"time"
)

func GuessHandler(code string) (map[string]interface{}, error) {
	bundleSize, metricMap, err := AstParser(string(code))
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	conn := Connect("guess.sock")
	defer conn.Close()
	bundleCode, err := BundleJSCode(code)
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return nil, err
	}
	b, _ := json.Marshal(map[string]interface{}{
		"bundleSize": bundleSize,
		"metrics":    metricMap,
		"code":       string(bundleCode),
	})
	_, err = conn.Write(b)
	if err != nil {
		return nil, err
	}

	// Read response
	var result []byte
	buf := make([]byte, 4096)

	for {
		conn.SetReadDeadline(time.Now().Add(5000 * time.Millisecond))
		n, err := conn.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
		}
		if err != nil {
			if err == io.EOF || len(result) > 0 {
				break
			}
			return nil, err
		}
	}

	var resultMap map[string]interface{}

	err = json.Unmarshal(result, &resultMap)
	if err != nil {
		return nil, err
	}

	resultMap["bundleSize"] = bundleSize
	resultMap["metrics"] = metricMap

	return resultMap, nil
}
