package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	Description string  `json:"description"`
	IconURL     string  `json:"icon_url"`
	Error       string  `json:"error,omitempty"`
}

func getWeather(city string) (*WeatherResponse, error) {
	// Memuat variabel lingkungan dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	response, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + city)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching data")
	}

	var data map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Mendapatkan data dari respons API
	current := data["current"].(map[string]interface{})
	location := data["location"].(map[string]interface{})
	condition := current["condition"].(map[string]interface{})
	iconURL := "http:" + condition["icon"].(string)

	return &WeatherResponse{
		City:        location["name"].(string),
		Temperature: current["temp_c"].(float64),
		Humidity:    int(current["humidity"].(float64)),
		Description: condition["text"].(string),
		IconURL:     iconURL,
	}, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Menambahkan CORS agar semua asal bisa mengakses
	w.Header().Set("Content-Type", "application/json")

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}

	weather, err := getWeather(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func main() {
	http.HandleFunc("/weather", weatherHandler)
	fmt.Println("Weather service is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
