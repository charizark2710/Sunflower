package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rogchap.com/v8go"

	LogConstant "RDIPs-Task/constant/LogConst"
	redis "RDIPs-Task/handler/Redis"
	unixSock "RDIPs-Task/handler/UnixSock"
	"RDIPs-Task/utils"
)

func ExecutionHandler(id string, msg map[string]interface{}) (any, error) {
	conn := unixSock.Connect("guess.sock")
	defer conn.Close()

	code := msg["code"].(string)
	ic := msg["ic"].(float64)
	cycle := msg["cycle"].(float64)
	input_ids := msg["input_ids"]
	attention_mask := msg["attention_mask"]
	// Send request to model
	b, _ := json.Marshal(map[string]any{
		"pred_ic":        ic,
		"pred_cycle":     cycle,
		"input_ids":      input_ids,
		"attention_mask": attention_mask,
		"ast_metric":     msg["ast_metric"],
	})
	b = append(b, []byte("#END#")...)

	for i := 0; i < len(b); i += 4096 {
		end := min(i+4096, len(b))
		if _, err := conn.Write(b[i:end]); err != nil {
			return nil, err
		}
	}
	done := false
	var result *v8go.Value
	for !done {
		// Read response
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		guessResult := make(map[string]any)
		if err := json.Unmarshal(buf[:n], &guessResult); err != nil {
			return nil, err
		}

		// Safe type extraction
		isSkip, _ := guessResult["isSkip"].(bool)
		done, _ = guessResult["done"].(bool)
		confidence, _ := guessResult["confidence"].(float64)

		// Decision logic
		if isSkip {
			utils.Log(LogConstant.Warning, "Skipping execution")
			continue
		}

		// Use separate contexts
		execCtx, execCancel := context.WithTimeout(context.Background(), 5*time.Minute)
		releaseCtx, releaseCancel := context.WithTimeout(context.Background(), 10*time.Second)

		redisLocker := redis.GetRedisLockerClient()
		lock, err := redisLocker.Obtain(execCtx, id, time.Minute*5, nil)
		if err != nil {
			execCancel()
			releaseCancel()
			return nil, err
		}

		resultCh := make(chan map[string]any, 1)

		go func() {
			val, actual_ic, actual_cycle, err := ExecuteJs(code, confidence < 0.7)
			// Send only if context not cancelled
			select {
			case resultCh <- map[string]any{
				"result": val,
				"ic":     actual_ic,
				"cycle":  actual_cycle,
				"err":    err,
			}:
			case <-execCtx.Done():
				return
			}
		}()

		var actual_ic float64
		var actual_cycle float64

		select {
		case <-execCtx.Done():
			lock.Release(releaseCtx)
			execCancel()
			releaseCancel()
			return nil, fmt.Errorf("execution timed out after 5 minutes")
		case r := <-resultCh:
			if r["err"] != nil {
				lock.Release(releaseCtx)
				execCancel()
				releaseCancel()
				return nil, r["err"].(error)
			}
			result = r["result"].(*v8go.Value)
			actual_ic = r["ic"].(float64)
			actual_cycle = r["cycle"].(float64)
		}

		if confidence < 0.6 {
			// send actual ic to model
			b, _ := json.Marshal(map[string]any{
				"actual_ic":    actual_ic,
				"actual_cycle": actual_cycle,
			})
			b = append(b, []byte("#END#")...)
			if _, err := conn.Write(b); err != nil {
				lock.Release(releaseCtx)
				execCancel()
				releaseCancel()
				return nil, err
			}
		} else {
			msg, _ := json.Marshal(map[string]any{
				"status": "complete",
			})
			msg = append(msg, []byte("#END#")...)
			if _, err := conn.Write(msg); err != nil {
				lock.Release(releaseCtx)
				execCancel()
				releaseCancel()
				return nil, err
			}
		}

		lock.Release(releaseCtx)
		execCancel()
		releaseCancel()
	}

	if result != nil {
		return result.String(), nil
	}
	return nil, nil
}
