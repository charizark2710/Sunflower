package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"rdips-dashboard/internal/config"
	"rdips-dashboard/internal/models"
)

type APIService struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewAPIService(cfg config.Config) *APIService {
	return &APIService{
		baseURL: cfg.APIBaseURL,
		apiKey:  cfg.APIKey,
		client:  &http.Client{Timeout: cfg.HTTPTimeout},
	}
}

func (s *APIService) Upstream(ctx context.Context, method, upstreamPath string, body io.Reader, contentType string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+upstreamPath, body)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("X-API-Key", s.apiKey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return resp.StatusCode, data, err
}

func (s *APIService) FetchDevices(ctx context.Context) ([]models.Device, error) {
	status, body, err := s.Upstream(ctx, http.MethodGet, "/devices", nil, "")
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("failed to fetch devices (%d)", status)
	}
	var envelope models.APIEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}
	var devices []models.Device
	if len(envelope.Data) > 0 && string(envelope.Data) != "null" {
		if err := json.Unmarshal(envelope.Data, &devices); err != nil {
			return nil, err
		}
	}
	return devices, nil
}

func (s *APIService) FetchDeviceDetail(ctx context.Context, id string, detail bool) (models.Device, error) {
	upstreamPath := "/devices/" + url.PathEscape(id)
	if detail {
		upstreamPath += "?detail=true"
	}
	status, body, err := s.Upstream(ctx, http.MethodGet, upstreamPath, nil, "")
	if err != nil {
		return models.Device{}, err
	}
	if status < 200 || status >= 300 {
		return models.Device{}, fmt.Errorf("failed to fetch device detail (%d)", status)
	}
	var envelope models.APIEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return models.Device{}, err
	}
	var device models.Device
	if len(envelope.Data) > 0 && string(envelope.Data) != "null" {
		if err := json.Unmarshal(envelope.Data, &device); err != nil {
			return models.Device{}, err
		}
	}
	return device, nil
}

func (s *APIService) FetchPerformancePoints(ctx context.Context, id string) ([]models.ChartPoint, error) {
	status, body, err := s.Upstream(ctx, http.MethodGet, "/performances/"+url.PathEscape(id), nil, "")
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("failed to fetch performance data (%d)", status)
	}
	var root map[string]any
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, err
	}
	data, _ := root["data"].(map[string]any)
	payload, _ := data["Payload"].([]any)
	if len(payload) == 0 {
		payload, _ = data["payload"].([]any)
	}
	points := make([]models.ChartPoint, 0, len(payload))
	for _, item := range payload {
		record, ok := item.(map[string]any)
		if !ok {
			continue
		}
		t := parseRecordTime(record)
		points = append(points, models.ChartPoint{
			Label:       t.Format("15:04"),
			LabelFull:   t.Format("2006-01-02 15:04:05"),
			Capacity:    parseNumber(record["capacity"]),
			InCapacity:  parseNumber(record["inCapacity"]),
			MaxCapacity: parseNumber(record["maxCapacity"]),
			OutCapacity: parseNumber(record["outCapacity"]),
		})
	}
	return points, nil
}
