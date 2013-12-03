package minimalweather

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var c = Pool.Get()

type CityWeather struct {
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
	Weather     Weather     `json:"weather"`
}

func outputWeatherAsJSON(current_city City, current_weather Weather) []byte {
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

	current_city := <-FindByName(city_name)
	current_weather := <-GetWeather(current_city.Coords)

	if current_city.Error != nil {
		http.NotFound(w, req)
	} else {
		out := outputWeatherAsJSON(current_city, current_weather)
		w.Write(out)
	}
}

func weatherByCoords(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lng, _ := strconv.ParseFloat(vars["lng"], 64)

	log.Println("By Coords:", lat, lng)

	coords := Coordinates{lat, lng}
	city_chan := FindByCoords(coords)
	current_weather := <-GetWeather(coords)
	current_city := <-city_chan

	out := outputWeatherAsJSON(current_city, current_weather)

	w.Write(out)
}

func Handler() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/weather/{city}", weatherByCity).Methods("GET")
	r.HandleFunc("/weather/{lat}/{lng}", weatherByCoords).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./website/")))

	return r
}
