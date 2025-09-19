package handler

import (
	LogConstant "RDIPs-Scheduler/constant/LogConst"
	"RDIPs-Scheduler/utils"
	"encoding/json"
	"io"
	"time"
)

func GuessHandler(code string) (string, error) {
	bundleSize, metricMap, err := AstParser(string(code))
	if err != nil {
		utils.Log(LogConstant.Error, err)
		return "", err
	}
	conn := Connect("/guess.sock")
	defer conn.Close()

	b, _ := json.Marshal(map[string]interface{}{
		"bundleSize": bundleSize,
		"metrics":    metricMap,
		"code":       code,
	})
	_, err = conn.Write(b)
	if err != nil {
		return "", err
	}

	// Read response
	var result []byte
	buf := make([]byte, 4096)

	for {
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, err := conn.Read(buf)
		if n > 0 {
			result = append(result, buf[:n]...)
		}
		if err != nil {
			if err == io.EOF || len(result) > 0 {
				break
			}
			return "", err
		}
	}

	return string(result), nil
}
