package models

import "encoding/json"

type APIEnvelope struct {
	HTTPCode int             `json:"httpCode"`
	Data     json.RawMessage `json:"data"`
	Message  string          `json:"message"`
	Error    any             `json:"Error"`
}

type Device struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Region        string `json:"region"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PerformanceID string `json:"performanceID"`
}

type ChartPoint struct {
	Label       string
	LabelFull   string
	Capacity    float64
	InCapacity  float64
	MaxCapacity float64
	OutCapacity float64
}

type ChartSeries struct {
	Name   string
	Color  string
	Points string
}

type AxisLabel struct {
	Text string
	X    float64
	Y    float64
}

type ChartView struct {
	ID       string
	Points   []ChartPoint
	Series   []ChartSeries
	MinY     float64
	MaxY     float64
	Ticks    []float64
	HasData  bool
	Message  string
	Width    float64
	Height   float64
	PlotX    float64
	PlotY    float64
	PlotW    float64
	PlotH    float64
	XLabels  []AxisLabel
	YLabels  []AxisLabel
	LastSeen string
}

func (c ChartView) X(i int) float64 {
	if len(c.Points) == 1 {
		return c.PlotX + c.PlotW/2
	}
	return c.PlotX + (float64(i)/float64(len(c.Points)-1))*c.PlotW
}

func (c ChartView) Y(value float64) float64 {
	return c.PlotY + c.PlotH - ((value-c.MinY)/(c.MaxY-c.MinY))*c.PlotH
}

type ScenarioRequest struct {
	Type            string  `json:"type"`
	DurationHours   float64 `json:"durationHours"`
	IntervalMinutes float64 `json:"intervalMinutes"`
}

type PayloadRecord struct {
	CreatedAt   string `json:"createdAt"`
	MaxCapacity string `json:"maxCapacity"`
	InCapacity  string `json:"inCapacity"`
	OutCapacity string `json:"outCapacity"`
	Capacity    string `json:"capacity"`
	UOM         string `json:"uom"`
}
