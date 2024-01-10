package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: weather [city name]")
		os.Exit(1)
	}

	city := os.Args[1]
	weatherData, err := getWeatherData(city)
	if err != nil {
		log.Fatalf("Error getting weather data: %v", err)
	}

	displayWeather(weatherData)
}

// WeatherData represents the top-level structure of the JSON response
type WeatherData struct {
	Location LocationData `json:"location"`
	Current  CurrentData  `json:"current"`
}

// LocationData represents the location part of the JSON
type LocationData struct {
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

// CurrentData represents the current weather data part of the JSON
type CurrentData struct {
	TempC float64 `json:"temp_c"`
}

func getWeatherData(city string) (*WeatherData, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY not set")
	}
	// mount the url with the city and api key and this url http://api.weatherapi.com/v1/current.json?key={apikey}&q={location}&aqi=no
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data WeatherData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func displayWeather(data *WeatherData) {
	fmt.Printf("Current Temperature in %s, %s: %.2f°C\n", data.Location.Name, data.Location.Region, data.Current.TempC)
}
