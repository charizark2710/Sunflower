package scenario

import (
	"crypto/rand"
	"math"
	mathrand "math/rand"
	"strconv"
	"time"

	"rdips-dashboard/internal/models"
)

func BuildPayload(scenarioType string, durationHours, intervalMinutes float64) []models.PayloadRecord {
	totalMinutes := durationHours * 60
	stepCount := max(1, int(math.Floor(totalMinutes/intervalMinutes)))
	records := make([]models.PayloadRecord, 0, stepCount+1)
	now := time.Now()
	rng := mathrand.New(mathrand.NewSource(randomSeed()))
	for i := 0; i <= stepCount; i++ {
		baseMax := randomInRange(rng, 8, 14, 2)
		maxCapacity := baseMax
		capacity := baseMax
		switch scenarioType {
		case "UNDERESTIMATE":
			capacity = randomInRange(rng, maxCapacity*0.3, maxCapacity*0.8, 2)
		case "OVERESTIMATE":
			capacity = randomInRange(rng, maxCapacity*1.2, maxCapacity*1.8, 2)
		}
		records = append(records, models.PayloadRecord{
			CreatedAt:   now.Add(time.Duration(i) * time.Duration(intervalMinutes) * time.Minute).Format("2/1/2006 15:04:05"),
			MaxCapacity: trimFloat(maxCapacity),
			InCapacity:  trimFloat(randomInRange(rng, capacity*0.1, capacity*0.4, 2)),
			OutCapacity: trimFloat(randomInRange(rng, capacity*0.02, capacity*0.15, 2)),
			Capacity:    trimFloat(capacity),
			UOM:         "kw",
		})
	}
	return records
}

func randomInRange(rng *mathrand.Rand, minValue, maxValue float64, decimals int) float64 {
	value := rng.Float64()*(maxValue-minValue) + minValue
	factor := math.Pow10(decimals)
	return math.Round(value*factor) / factor
}

func randomSeed() int64 {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().UnixNano()
	}
	var seed int64
	for _, value := range b {
		seed = seed<<8 + int64(value)
	}
	return seed
}

func trimFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
