package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bmizerany/pat"
	"github.com/elcuervo/minimalweather/city"
	"github.com/elcuervo/minimalweather/weather"
)

type CityWeather struct {
	Name        string           `json:"name"`
	Coordinates city.Coordinates `json:"coordinates"`
	Weather     weather.Weather  `json:"weather"`
}

func outputWeatherAsJSON(current_city city.City, current_weather weather.Weather) []byte {
	city_weather := &CityWeather{
		Name:        current_city.Name,
		Coordinates: current_city.Coords,
		Weather:     current_weather}

	out, _ := json.Marshal(city_weather)

	return out
}

func weatherByCity(w http.ResponseWriter, req *http.Request) {
	log.Println("By Name")

	city_name := req.URL.Query().Get(":city")

	current_city := <-city.FindByName(city_name)
	current_weather := <-weather.GetWeather(current_city.Coords)

	out := outputWeatherAsJSON(current_city, current_weather)

	w.Write(out)
}

func weatherByCoords(w http.ResponseWriter, req *http.Request) {
	log.Println("By Coords")

	lat, _ := strconv.ParseFloat(req.URL.Query().Get(":lat"), 64)
	lng, _ := strconv.ParseFloat(req.URL.Query().Get(":lng"), 64)

	coords := city.Coordinates{lat, lng}
	city_chan := city.FindByCoords(coords)
	current_weather := <-weather.GetWeather(coords)
	current_city := <-city_chan

	out := outputWeatherAsJSON(current_city, current_weather)

	w.Write(out)
}

func main() {
	m := pat.New()
	port := ":" + os.Getenv("PORT")

	m.Get("/weather/:city", http.HandlerFunc(weatherByCity))
	m.Get("/weather/:lat/:lng", http.HandlerFunc(weatherByCoords))
	http.Handle("/", m)

	log.Println("Listening in", port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
