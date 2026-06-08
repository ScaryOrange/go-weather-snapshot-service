package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type openMeteoClient struct {
	httpClient   *http.Client
	urlGeocoding string
	urlForecast  string
}

type geocodingResp struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

type weatherResp struct {
	CurrentWeather struct {
		Time        string  `json:"time"`
		Temperature float64 `json:"temperature"`
		WindSpeed   float64 `json:"windspeed"`
	} `json:"current_weather"`
}

func NewOpenMeteoClient(geoEndpoint string, forecastEndpoint string) *openMeteoClient {
	return &openMeteoClient{
		httpClient:   &http.Client{Timeout: 2 * time.Second},
		urlGeocoding: geoEndpoint,
		urlForecast:  forecastEndpoint,
	}
}

func (omc *openMeteoClient) doGet(ctx context.Context, myUrl string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", myUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := omc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (omc *openMeteoClient) fromCityToCoordinates(ctx context.Context, city string) (float64, float64, error) {
	params := url.Values{}
	params.Set("name", city)
	params.Set("count", "1")
	params.Set("language", "en")
	params.Set("format", "json")
	myUrl := omc.urlGeocoding + "?" + params.Encode()

	var geoResp geocodingResp

	body, err := omc.doGet(ctx, myUrl)
	if err != nil {
		return 0, 0, err
	}

	err = json.Unmarshal(body, &geoResp)
	if err != nil {
		return 0, 0, err
	}

	if len(geoResp.Results) == 0 {
		return 0, 0, fmt.Errorf("City not found")
	}

	return geoResp.Results[0].Latitude, geoResp.Results[0].Longitude, nil
}

func (omc *openMeteoClient) fetchWeather(ctx context.Context, latitude float64, longitude float64) (weatherResp, []byte, error) {
	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%f", latitude))
	params.Set("longitude", fmt.Sprintf("%f", longitude))
	params.Set("current_weather", "true")
	myUrl := omc.urlForecast + "?" + params.Encode()

	var resp weatherResp

	body, err := omc.doGet(ctx, myUrl)
	if err != nil {
		return weatherResp{}, nil, err
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return weatherResp{}, nil, err
	}

	return resp, body, nil
}

func (omc *openMeteoClient) Current(ctx context.Context, city string) (WeatherSnapshot, error) {
	if len(city) < 2 {
		return WeatherSnapshot{}, fmt.Errorf("City not found")
	}

	latitude, longitude, err := omc.fromCityToCoordinates(ctx, city)
	if err != nil {
		return WeatherSnapshot{}, err
	}

	response, rawResponse, err := omc.fetchWeather(ctx, latitude, longitude)
	if err != nil {
		return WeatherSnapshot{}, err
	}

	snapshot := WeatherSnapshot{
		City:               CityNormalize(city),
		Provider:           "open-meteo",
		TemperatureCelsius: response.CurrentWeather.Temperature,
		WindSpeed:          response.CurrentWeather.WindSpeed,
		ObservedAt:         response.CurrentWeather.Time,
		RawPayload:         rawResponse,
	}

	return snapshot, nil
}
