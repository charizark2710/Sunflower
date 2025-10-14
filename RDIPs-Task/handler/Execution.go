package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"rogchap.com/v8go"

	redis "RDIPs-Task/handler/Redis"
	unixSock "RDIPs-Task/handler/UnixSock"
)

func ExecutionHandler(id string, msg map[string]interface{}) (any, error) {
	conn := unixSock.Connect("/guess.sock")
	defer conn.Close()

	code := msg["code"].(string)
	ic := msg["ic"].(float64)
	cycle := msg["cycle"].(float64)
	input_ids := msg["input_ids"].([]float64)
	attention_mask := msg["attention_mask"].([]float64)
	// Send request to model
	b, _ := json.Marshal(map[string]any{
		"pred_ic":        ic,
		"pred_cycle":     cycle,
		"input_ids":      input_ids,
		"attention_mask": attention_mask,
	})
	b = append(b, []byte("#END#\n")...)

	for i := 0; i < len(b); i += 4096 {
		if _, err := conn.Write(b[i:4096]); err != nil {
			return nil, err
		}
	}

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
	isSkip, _ := guessResult["isSkip"].(string)
	confidence, _ := guessResult["confidence"].(float64)

	// Decision logic
	if isSkip == "true" {
		return nil, fmt.Errorf("Skip")
	}

	// Execute if not skipped
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	redisLocker := redis.GetRedisLockerClient()
	lock, err := redisLocker.Obtain(ctx, id, time.Minute*5, nil)

	if err != nil {
		return nil, err
	}

	defer lock.Release(ctx)

	resultCh := make(chan map[string]any, 1)

	go func() {
		val, actual_ic, actual_cycle, err := ExecuteJs(code, confidence < 0.7)
		resultCh <- map[string]any{
			"result": val,
			"ic":     actual_ic,
			"cycle":  actual_cycle,
			"err":    err,
		}
	}()

	var result *v8go.Value
	var actual_ic float64
	var actual_cycle float64
	select {
	case <-ctx.Done():
		if _, err := conn.Write(nil); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("execution timed out after 5 minutes")
	case r := <-resultCh:
		if r["err"] != nil {
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
		b = append(b, []byte("#END#\n")...)
		if _, err := conn.Write(b); err != nil {
			return nil, err
		}
	}

	// Try parse as JSON if possible
	var out any
	if json.Unmarshal([]byte(result.String()), &out) == nil {
		return out, nil
	}

	// Fallback to raw string
	return result.String(), nil
}
