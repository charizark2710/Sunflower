package middleware

import (
	"RDIPs-BE/model"
	"bytes"
	"context"
	"encoding/json"
	"time"

	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	WeatherAuthenKey string = ""
	//Set timeout for authetication key of Weather API Login
	refreshPeriod   = 1 * time.Minute
	lastFetchedTime = time.Now()
)

func ValidationAPIWeatherKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		if keyExpired() {
			ctx, cancel := context.WithTimeout(c.Request.Context(), refreshPeriod)
			defer cancel()

			err := callWeatherLogin(ctx, c)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		c.Next()
	}
}

func callWeatherLogin(ctx context.Context, c *gin.Context) error {
	weatherAccount := model.WeatherAccount{
		UserName: os.Getenv("WEATHER_USER_NAME"),
		Password: os.Getenv("WEATHER_PASSWORD"),
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(weatherAccount)
	if err == nil {
		url := os.Getenv("WEATHER_BASE_URL") + "/bigdata-weather/v1/guest/login"
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("apikey", os.Getenv("WEATHER_API_KEY_HEADER"))

		client := &http.Client{
			Timeout: 5 * time.Second, // Set a timeout for the API call
		}

		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var body, _ = io.ReadAll(resp.Body)
			var loginResponse model.LoginResponse
			if err := json.Unmarshal(body, &loginResponse); err == nil {
				WeatherAuthenKey = loginResponse.Result.ApiKey
				lastFetchedTime = time.Now()
			}
		}
	}
	return err

}

func keyExpired() bool {
	if WeatherAuthenKey == "" {
		return true
	}

	return time.Now().After(lastFetchedTime.Add(refreshPeriod))
}
