package chart

import (
	"fmt"
	"math"
	"strings"

	"rdips-dashboard/internal/models"
)

func Build(id string, points []models.ChartPoint) models.ChartView {
	view := models.ChartView{
		ID:      id,
		Points:  points,
		HasData: len(points) > 0,
		Width:   960,
		Height:  430,
		PlotX:   58,
		PlotY:   24,
		PlotW:   860,
		PlotH:   320,
		Message: "No data available",
	}
	if len(points) == 0 {
		return view
	}
	view.LastSeen = points[len(points)-1].LabelFull
	values := make([]float64, 0, len(points)*4)
	for _, p := range points {
		values = append(values, p.Capacity, p.InCapacity, p.MaxCapacity, p.OutCapacity)
	}
	minY, maxY := values[0], values[0]
	for _, value := range values {
		minY = math.Min(minY, value)
		maxY = math.Max(maxY, value)
	}
	if minY == maxY {
		minY -= 1
		maxY += 1
	}
	padding := (maxY - minY) * 0.12
	view.MinY = math.Max(0, minY-padding)
	view.MaxY = maxY + padding
	for i := 0; i < 5; i++ {
		tick := view.MinY + (view.MaxY-view.MinY)*float64(i)/4
		view.Ticks = append(view.Ticks, tick)
		y := view.PlotY + view.PlotH - ((tick-view.MinY)/(view.MaxY-view.MinY))*view.PlotH
		view.YLabels = append(view.YLabels, models.AxisLabel{Text: fmt.Sprintf("%.1f", tick), X: view.PlotX - 12, Y: y + 4})
	}
	for i, p := range points {
		if len(points) > 8 && i%(int(math.Ceil(float64(len(points))/8))) != 0 && i != len(points)-1 {
			continue
		}
		x := view.X(i)
		view.XLabels = append(view.XLabels, models.AxisLabel{Text: p.Label, X: x, Y: view.PlotY + view.PlotH + 34})
	}
	series := []struct {
		name  string
		color string
		value func(models.ChartPoint) float64
	}{
		{"Capacity (kW)", "#2563eb", func(p models.ChartPoint) float64 { return p.Capacity }},
		{"In Capacity (kW)", "#16a34a", func(p models.ChartPoint) float64 { return p.InCapacity }},
		{"Max Capacity (kW)", "#d97706", func(p models.ChartPoint) float64 { return p.MaxCapacity }},
		{"Out Capacity (kW)", "#dc2626", func(p models.ChartPoint) float64 { return p.OutCapacity }},
	}
	for _, s := range series {
		coords := make([]string, 0, len(points))
		for i, p := range points {
			coords = append(coords, fmt.Sprintf("%.1f,%.1f", view.X(i), view.Y(s.value(p))))
		}
		view.Series = append(view.Series, models.ChartSeries{Name: s.name, Color: s.color, Points: strings.Join(coords, " ")})
	}
	return view
}
