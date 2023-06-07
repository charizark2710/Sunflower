package model

import (
	"encoding/json"
	"strings"
	"time"
)

type WeatherAccount struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status int            `json:"status"`
	Result ApiKeyResponse `json:"result"`
}

type ApiKeyResponse struct {
	ApiKey string `json:"ApiKey"`
}

type WeatherRequest struct {
	FromDate string `json:"fromDate,omitempty"`
	ToDate   string `json:"toDate"`
	Latlongs string `json:"latlongs"`
	Period   int    `json:"period"`
}

type WeatherResponse struct {
	Status int           `json:"status"`
	Result []WeatherData `json:"result"`
}

type WeatherData struct {
	Time time.Time   `json:"time"`
	Info WeatherInfo `json:"info"`
}

type WeatherInfo struct {
	WeatherSummary       string  `json:"weather_summary_en"`
	FeelsLikeTemperature float64 `json:"feels_like_temperature"`
	Precipitation        float64 `json:"precipitation"`
	ProbabilityRain      float64 `json:"probability_rain"`
	AirTemperature       float64 `json:"air_temperature"`
	AirTemperatureMin    float64 `json:"air_temperature_min"`
	AirTemperatureMax    float64 `json:"air_temperature_max"`
	CloudCover           float64 `json:"cloud_cover"`
	Weather              float64 `json:"weather"`
	WeatherStatus        float64 `json:"weather_status"`
	WindDirection        float64 `json:"wind_direction"`
	WindGust             float64 `json:"wind_gust"`
	WindSpeed            float64 `json:"wind_speed"`
	WindDirectionEn      string  `json:"wind_direction_en"`
	RelativeHumidity     float64 `json:"relative_humidity"`
	DewPoint             float64 `json:"dew_point"`
	Radiation            float64 `json:"radiation"`
	Sunrise              int     `json:"sunrise"`
	SunSet               int     `json:"sunset"`
}

// Custom unmarshal function to map "GMOS" fields to "Info" in the struct
func (r *WeatherData) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	for key, val := range v {
		if strings.HasPrefix(key, "GMOS") {
			info := WeatherInfo{}
			if resp, err := json.Marshal(val); err == nil {
				if err := json.Unmarshal(resp, &info); err != nil {
					return err
				}
				r.Info = info
			} else {
				return err
			}
		}

		if key == "time" {
			//Convert string to time.Time
			timeStr := val.(string)
			layout := "15:04 02/01/2006"

			timeObj, err := time.Parse(layout, timeStr)
			if err != nil {
				return err
			}
			r.Time = timeObj
			// r.Time = val.(string)
		}
	}
	return nil
}
