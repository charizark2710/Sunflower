package services

import (
	LogConstant "RDIPs-BE/constant/LogConst"
	"RDIPs-BE/model"
	commonModel "RDIPs-BE/model/common"
	"RDIPs-BE/utils"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

var GetWeatherNext14Days = func(c *commonModel.ServiceContext) (commonModel.ResponseTemplate, error) {
	utils.Log(LogConstant.Info, "GetWeatherNext14Days Start")
	if authen := GetAuthorization(); authen != "" {
		weatherBody := model.WeatherRequest{}
		if err := json.Unmarshal(c.Body, &weatherBody); err == nil {

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(weatherBody); err == nil {
				url := os.Getenv("WEATHER_BASE_URL") + "14days"
				req, _ := http.NewRequest("POST", url, &buf)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Add("apikey", os.Getenv("WEATHER_API_KEY_HEADER"))
				req.Header.Set("Authorization", "Bearer "+authen)

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					var body, _ = io.ReadAll(resp.Body)
					var weatherForecastResponse model.WeatherResponse
					if err := json.Unmarshal(body, &weatherForecastResponse); err == nil {
						return commonModel.ResponseTemplate{HttpCode: 200, Data: weatherForecastResponse}, nil
					}
				}
			}
			utils.Log(LogConstant.Info, "GetWeatherNext14Days End")
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		} else {
			return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, err
		}
	} else {
		return commonModel.ResponseTemplate{HttpCode: 500, Data: nil}, nil
	}

}

var GetAuthorization = func() string {
	utils.Log(LogConstant.Info, "GetAuthorization Start")
	var apikey = ""
	weatherAccount := model.WeatherAccount{
		UserName: os.Getenv("WEATHER_USER_NAME"),
		Password: os.Getenv("WEATHER_PASSWORD"),
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(weatherAccount); err == nil {
		req, _ := http.NewRequest("POST", os.Getenv("WEATHER_LOGIN_URL"), &buf)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("apikey", os.Getenv("WEATHER_API_KEY_HEADER"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			utils.Log(LogConstant.Info, "GetAuthorization End")
			var body, _ = io.ReadAll(resp.Body)
			var loginResponse model.LoginResponse
			if err := json.Unmarshal(body, &loginResponse); err == nil {
				apikey = loginResponse.Result.ApiKey
			}
		}
	}

	return apikey

}
