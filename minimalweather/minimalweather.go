package minimalweather

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/elcuervo/geoip"
        "github.com/ianoshen/uaparser"
	"log"
	"math"
	"os"
	"net/http"
        "html/template"
	"strconv"
)

var c = Pool.Get()

type CityWeather struct {
	Name        string      `json:"name"`
	Coordinates Coordinates `json:"coordinates"`
	Weather     Weather     `json:"weather"`
        JSON        string      `json:"-"`
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

func iOSDevice(agent string) bool {
        u := uaparser.Parse(agent)
        return u.Device.Name == "iPhone" || u.Device.Name == "iPad" || u.Device.Name == "iPod"
}

func geolocate(req *http.Request) geoip.Geolocation {
        var user_addr string
        var current_env = os.Getenv("DEVELOPMENT")

        if current_env != "" {
                user_addr = "186.54.253.75"
        } else {
                user_addr = req.RemoteAddr
        }

        return <-GetLocation(user_addr)
}

func homepage(w http.ResponseWriter, req *http.Request) {
        var cw *CityWeather

 //       if iOSDevice(req.UserAgent()) {
                geo := geolocate(req)
                coords := Coordinates{geo.Location.Latitude, geo.Location.Longitude}
                city := <-FindByCoords(coords)
                weather := <-GetWeather(coords)

                cw = &CityWeather{
                        Name: city.Name,
                        Coordinates: city.Coords,
                        Weather: weather,
                }

//        }

        t, _ := template.ParseFiles("./website/index.html")
        out, err := json.Marshal(cw)
        cw.JSON = string(out)
        cw.Weather.Temperature = math.Floor(cw.Weather.Temperature)
        err = t.Execute(w, cw)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

func Handler() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/weather/{city}", weatherByCity).Methods("GET")
	r.HandleFunc("/weather/{lat}/{lng}", weatherByCoords).Methods("GET")

	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./website/")))

	r.HandleFunc("/", homepage).Methods("GET")

	return r
}
