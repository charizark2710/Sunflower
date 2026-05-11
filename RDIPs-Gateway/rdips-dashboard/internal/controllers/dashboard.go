package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"rdips-dashboard/internal/chart"
	"rdips-dashboard/internal/models"
	"rdips-dashboard/internal/scenario"
	"rdips-dashboard/internal/services"
	"rdips-dashboard/internal/views"
)

type DashboardController struct {
	api      *services.APIService
	renderer *views.Renderer
}

func NewDashboardController(api *services.APIService, renderer *views.Renderer) *DashboardController {
	return &DashboardController{api: api, renderer: renderer}
}

func (c *DashboardController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", c.Home)
	mux.HandleFunc("GET /partials/device-selector", c.DeviceSelector)
	mux.HandleFunc("POST /select-device", c.SelectDevice)
	mux.HandleFunc("GET /partials/performance-chart/", c.PerformanceChart)
	mux.HandleFunc("GET /api/devices", c.ProxyDevices)
	mux.HandleFunc("GET /api/devices/", c.ProxyDeviceDetail)
	mux.HandleFunc("POST /api/devices/", c.PostDeviceScenario)
	mux.HandleFunc("GET /api/performances/", c.ProxyPerformance)
}

func (c *DashboardController) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		c.Dashboard(w, r)
		return
	}
	c.renderer.HTML(w, "home.html", map[string]any{"Title": "RDIPs Dashboard"})
}

func (c *DashboardController) Dashboard(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[1] != "dashboard" || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	c.renderer.HTML(w, "dashboard.html", map[string]any{
		"Title": "Device Dashboard",
		"ID":    parts[0],
	})
}

func (c *DashboardController) DeviceSelector(w http.ResponseWriter, r *http.Request) {
	devices, err := c.api.FetchDevices(r.Context())
	data := map[string]any{
		"Devices": devices,
		"Error":   "",
	}
	if err != nil {
		data["Error"] = err.Error()
	}
	c.renderer.HTML(w, "device_selector.html", data)
}

func (c *DashboardController) SelectDevice(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	id := strings.TrimSpace(r.FormValue("device_id"))
	if id == "" {
		c.renderer.HTML(w, "device_selector.html", map[string]any{
			"Devices": []models.Device{},
			"Error":   "Choose a device first.",
		})
		return
	}
	detail, err := c.api.FetchDeviceDetail(r.Context(), id, true)
	if err != nil {
		devices, _ := c.api.FetchDevices(r.Context())
		c.renderer.HTML(w, "device_selector.html", map[string]any{
			"Devices": devices,
			"Error":   err.Error(),
		})
		return
	}
	if detail.PerformanceID == "" {
		devices, _ := c.api.FetchDevices(r.Context())
		c.renderer.HTML(w, "device_selector.html", map[string]any{
			"Devices": devices,
			"Error":   "This device has no performance ID.",
		})
		return
	}
	w.Header().Set("HX-Redirect", "/"+url.PathEscape(detail.PerformanceID)+"/dashboard")
	w.WriteHeader(http.StatusNoContent)
}

func (c *DashboardController) PerformanceChart(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/partials/performance-chart/")
	id = path.Clean("/" + id)[1:]
	if id == "." || id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	points, err := c.api.FetchPerformancePoints(r.Context(), id)
	view := chart.Build(id, points)
	if err != nil {
		view.Message = err.Error()
		view.HasData = false
	}
	c.renderer.HTML(w, "performance_chart.html", view)
}

func (c *DashboardController) ProxyDevices(w http.ResponseWriter, r *http.Request) {
	c.proxy(w, r, http.MethodGet, "/devices", nil)
}

func (c *DashboardController) ProxyDeviceDetail(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	if strings.HasSuffix(rest, "/scenario") {
		http.NotFound(w, r)
		return
	}
	upstream := "/devices/" + rest
	if r.URL.RawQuery != "" {
		upstream += "?" + r.URL.RawQuery
	}
	c.proxy(w, r, http.MethodGet, upstream, nil)
}

func (c *DashboardController) ProxyPerformance(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/performances/")
	c.proxy(w, r, http.MethodGet, "/performances/"+id, nil)
}

func (c *DashboardController) PostDeviceScenario(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	deviceID := strings.TrimSuffix(rest, "/scenario")
	if deviceID == rest || deviceID == "" {
		http.NotFound(w, r)
		return
	}

	var body models.ScenarioRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "Invalid JSON body"})
		return
	}
	if body.Type != "UNDERESTIMATE" && body.Type != "OVERESTIMATE" && body.Type != "NORMAL" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": `Invalid "type". Must be one of: UNDERESTIMATE, OVERESTIMATE, NORMAL.`,
		})
		return
	}
	durationHours := body.DurationHours
	if durationHours <= 0 {
		durationHours = 4
	}
	intervalMinutes := body.IntervalMinutes
	if intervalMinutes <= 0 {
		intervalMinutes = 60
	}

	detail, err := c.api.FetchDeviceDetail(r.Context(), deviceID, true)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "Failed to fetch device detail", "detail": err.Error()})
		return
	}
	if detail.ID == "" {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": "Device has no performanceID", "device": detail})
		return
	}

	payload := scenario.BuildPayload(body.Type, durationHours, intervalMinutes)
	putBody, _ := json.Marshal(map[string]any{
		"document_name": detail.Name,
		"payload":       payload,
	})
	status, responseBody, err := c.api.Upstream(r.Context(), http.MethodPut, "/performances/"+url.PathEscape(detail.ID), bytes.NewReader(putBody), "application/json")
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "Failed to update performance", "detail": err.Error()})
		return
	}
	var upstream any = string(responseBody)
	_ = json.Unmarshal(responseBody, &upstream)
	if status < 200 || status >= 300 {
		writeJSON(w, status, map[string]any{"error": "Failed to update performance", "upstreamStatus": status, "upstreamBody": upstream})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":              true,
		"scenario":        body.Type,
		"deviceId":        deviceID,
		"deviceID":        detail.ID,
		"documentName":    detail.Name,
		"durationHours":   durationHours,
		"intervalMinutes": intervalMinutes,
		"recordCount":     len(payload),
		"payload":         payload,
		"upstream":        upstream,
	})
}

func (c *DashboardController) proxy(w http.ResponseWriter, r *http.Request, method, upstreamPath string, body io.Reader) {
	status, data, err := c.api.Upstream(r.Context(), method, upstreamPath, body, "")
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": "Upstream request failed", "detail": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
