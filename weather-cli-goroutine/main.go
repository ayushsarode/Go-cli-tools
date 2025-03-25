package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"time"
	"github.com/joho/godotenv"
)

const BASE_URL = "https://api.openweathermap.org/data/2.5/weather"

type WeatherResponse struct {
	Name string	`json: "name"`
	Main struct {
		Temp float64 `json: "temp"`
	} `json:"main"`
	Weather []struct {
		Main string `json: "main"`
	}`json:"weather"`
}


func fetchWeather(city string, ch chan<- string){
	apikey := os.Getenv("WEATHER_API_KEY")

	if apikey == "" {
		fmt.Printf("API key is missing")
		return
	}
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", BASE_URL, city, apikey)

	
	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)

	if err != nil {
		ch <- fmt.Sprintf("Error fetching weather for %s: %v", city , err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ch <- fmt.Sprintf("City %s not found", city)
		return
	}

	var data WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		ch <- fmt.Sprintf("Error decoding resposne for %s: %v", city, err)
		return
	}

	ch <- fmt.Sprintf("%s: %.1f°C, %s", data.Name, data.Main.Temp, data.Weather[0].Main)



}

func main() {
	godotenv.Load()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [city1] [city2]...")
		return
	}


	cities := os.Args[1:]
	ch := make(chan string)

	for _, city := range cities {
		go fetchWeather(city, ch)
	}


	for range cities {
		fmt.Println(<-ch)
	}
	fmt.Println(len(os.Args))
}