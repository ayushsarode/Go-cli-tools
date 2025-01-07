package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const BASE_URL = "https://api.openweathermap.org/data/2.5/weather"

type WeatherResponse struct {
    Name string `json:"name"`
    Main struct {
        Temp      float64 `json:"temp"`
        FeelsLike float64 `json:"feels_like"`
        Humidity  float64 `json:"humidity"`
    } `json:"main"`
    Weather []struct {
        Description string `json:"description"`
    } `json:"weather"`
    Cod interface{} `json:"cod"`
    Msg string      `json:"message"`
}

func getCodAsInt(cod interface{}) int {
    switch value := cod.(type) {
    case float64:
        return int(value)
    case string:
        if num, err := strconv.Atoi(value); err == nil {
            return num
        }
    }
    return 0
}

func getWeather(city string) {
    apiKey := os.Getenv("WEATHER_API_KEY")
    if apiKey == "" {
        fmt.Println("❌ API key is missing. Set the WEATHER_API_KEY environment variable.")
        return
    }

    url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", BASE_URL, city, apiKey)

    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    var weatherData WeatherResponse
    if err := json.Unmarshal(body, &weatherData); err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    codInt := getCodAsInt(weatherData.Cod)
    if codInt != 200 {
        fmt.Println("❌ Invalid city name. Please check the spelling and try again.")
        return
    }

    fmt.Printf("Weather in %s:\n%s, %.2f°C\n", weatherData.Name, weatherData.Weather[0].Description, weatherData.Main.Temp)
    fmt.Printf("Feels like: %.2f°C \nHumidity: %.2f%%\n", weatherData.Main.FeelsLike, weatherData.Main.Humidity)
}

func main() {
    godotenv.Load() 

    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go <city>")
        return
    }
    city := os.Args[1]
    getWeather(city)
}
