package pages

import (
	"encoding/json"
	"net/http"
	"strconv"
	"log"
	"github.com/gorilla/mux"
	mw "github.com/elcuervo/minimalweather/minimalweather"
)

type API struct {}

func (api *API) CityByCoords(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lng, _ := strconv.ParseFloat(vars["lng"], 64)

	log.Println("city By Coords:", lat, lng)

	coords := mw.Coordinates{ lat, lng }
	current_city := <-mw.FindByCoords(coords)

        out, _ := json.Marshal(current_city)
	w.Write(out)
}

func (api *API) WeatherByCoords(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lng, _ := strconv.ParseFloat(vars["lng"], 64)

	log.Println("By Coords:", lat, lng)

	coords := mw.Coordinates{ lat, lng }
	current_city := mw.FindByCoords(coords)
	current_weather := <-mw.GetWeather(coords)

	out := api.outputWeatherAsJSON(<-current_city, current_weather)

	w.Write(out)
}

func (api *API) outputWeatherAsJSON(current_city mw.City, current_weather mw.Weather) []byte {
	city_weather := &CityWeather{
		City:        current_city,
		Weather:     current_weather}

	out, _ := json.Marshal(city_weather)

	return out
}

func (api *API) WeatherByCity(w http.ResponseWriter, req *http.Request) {
	city_name := mux.Vars(req)["city"]

	log.Println("By Name:", city_name)

	current_city := <-mw.FindByName(city_name)
	current_weather := <-mw.GetWeather(current_city.Coords)

	if current_city.Error != nil {
		http.NotFound(w, req)
	} else {
		out := api.outputWeatherAsJSON(current_city, current_weather)
		w.Write(out)
	}
}

