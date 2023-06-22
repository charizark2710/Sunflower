package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/middleware"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

var GetWeatherForecast = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetWeatherForecast Start")
	weatherBody := model.WeatherRequest{}
	err := json.Unmarshal(c.Body, &weatherBody)
	if err == nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(weatherBody); err == nil {
			url := os.Getenv("WEATHER_BASE_URL") + "/weather/v1/guest/14days"
			req, _ := http.NewRequest("POST", url, &buf)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("apikey", os.Getenv("WEATHER_API_KEY_HEADER"))
			req.Header.Set("Authorization", "Bearer "+middleware.WeatherAuthenKey)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				var body, _ = io.ReadAll(resp.Body)
				var weatherForecastResponse model.WeatherResponse
				if err := json.Unmarshal(body, &weatherForecastResponse); err == nil {
					utils.Log(LogConstant.Info, "GetWeatherForecast End")
					return commonModel.ResponseTemplate{HttpCode: 200, Data: weatherForecastResponse}, nil
				}
			}
		}
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
	}
	return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err

}
