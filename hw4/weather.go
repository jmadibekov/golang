// Fetches the real-time weather data for given cities
// and outputs them in json file
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const APIBase = "https://api.weatherapi.com/v1/current.json"
const APIKey = "8367e706aa5f43ea830101139211010"

var cities = [5]string{"Almaty", "Nur-Sultan", "London", "Moscow", "New-York"}

type Weather struct {
	City        string
	Temperature float64 // in celsius
}

type WeatherData struct {
	CurrentTime time.Time
	Cities      []Weather
}

func getCurrentWeather(city string) float64 {
	APIUrl := fmt.Sprintf("%v?key=%v&q=%v", APIBase, APIKey, city)
	log.Println("API url is", APIUrl)

	response, err := http.Get(APIUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	// checking if status is not 400
	if response.StatusCode == 400 {
		log.Fatalln("Received 400 error for city", city)
	}

	// reading the body from the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	current := result["current"].(map[string]interface{})
	log.Println("The weather is last updated on", current["last_updated"].(string))
	currentDegree := current["temp_c"].(float64)
	log.Println("The weather for", city, "is", currentDegree, "celsius")

	return currentDegree
}

func convertToJSON(weatherData WeatherData) {
	jsonData, err := json.MarshalIndent(weatherData, "", "   ")

	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile("weather.json", jsonData, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	weatherData := WeatherData{CurrentTime: time.Now(), Cities: make([]Weather, 0)}

	for _, city := range cities {
		weather := Weather{City: city, Temperature: getCurrentWeather(city)}
		weatherData.Cities = append(weatherData.Cities, weather)
	}

	convertToJSON(weatherData)
}
