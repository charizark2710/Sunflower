package services

import (
	"strconv"
	"time"
)

func parseRecordTime(record map[string]any) time.Time {
	if timestamp, ok := record["timestamp"].(map[string]any); ok {
		if date, ok := timestamp["$date"].(map[string]any); ok {
			if raw, ok := date["$numberLong"].(string); ok {
				if ms, err := strconv.ParseInt(raw, 10, 64); err == nil {
					return time.UnixMilli(ms)
				}
			}
		}
	}
	for _, key := range []string{"createdAt", "createdAt"} {
		raw, _ := record[key].(string)
		if raw == "" {
			continue
		}
		for _, layout := range []string{time.RFC3339Nano, "2/1/2006 15:04:05", "2006-01-02T15:04:05Z07:00"} {
			if t, err := time.Parse(layout, raw); err == nil {
				return t
			}
		}
	}
	return time.Now()
}

func parseNumber(v any) float64 {
	switch value := v.(type) {
	case float64:
		return value
	case string:
		n, _ := strconv.ParseFloat(value, 64)
		return n
	default:
		return 0
	}
}
