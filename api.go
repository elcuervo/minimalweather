package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/elcuervo/minimalweather/city"
	"github.com/elcuervo/minimalweather/weather"
	"github.com/gorilla/mux"
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
	city_name := mux.Vars(req)["city"]

	log.Println("By Name:", city_name)

	current_city := <-city.FindByName(city_name)
	current_weather := <-weather.GetWeather(current_city.Coords)

	out := outputWeatherAsJSON(current_city, current_weather)

	w.Write(out)
}

func weatherByCoords(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lng, _ := strconv.ParseFloat(vars["lng"], 64)

	log.Println("By Coords:", lat, lng)

	coords := city.Coordinates{lat, lng}
	city_chan := city.FindByCoords(coords)
	current_weather := <-weather.GetWeather(coords)
	current_city := <-city_chan

	out := outputWeatherAsJSON(current_city, current_weather)

	w.Write(out)
}

func main() {
	r := mux.NewRouter()
	port := ":" + os.Getenv("PORT")

	r.HandleFunc("/weather/{city}", weatherByCity).Methods("GET")
	r.HandleFunc("/weather/{lat}/{lng}", weatherByCoords).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./website/")))

	http.Handle("/", r)

	log.Println("Listening in", port)
	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
