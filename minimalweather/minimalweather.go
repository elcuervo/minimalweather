package minimalweather

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/elcuervo/geoip"
	"log"
	"os"
	"fmt"
	"math"
	"strings"
	"net/http"
        "html/template"
	"strconv"
)

var c = Pool.Get()

type CityWeather struct {
        City        City        `json:"city"`
	Weather     Weather     `json:"weather"`
        Unit        string      `json:"unit"`
        JSON        string      `json:"-"`
}

func outputWeatherAsJSON(current_city City, current_weather Weather) []byte {
	city_weather := &CityWeather{
		City:        current_city,
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

	coords := Coordinates{ lat, lng }
	current_city := FindByCoords(coords)
	current_weather := <-GetWeather(coords)

	out := outputWeatherAsJSON(<-current_city, current_weather)

	w.Write(out)
}

func cityByCoords(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lng, _ := strconv.ParseFloat(vars["lng"], 64)

	log.Println("city By Coords:", lat, lng)

	coords := Coordinates{ lat, lng }
	current_city := <-FindByCoords(coords)

        out, _ := json.Marshal(current_city)
	w.Write(out)
}

func ipFromRemote(s string) string {
        index := strings.LastIndex(s, ":")
        if index == -1 {
                return s
        }
        return s[:index]
}

func ipAddress(r *http.Request) string {
        development := os.Getenv("DEVELOPMENT")
        if development != "" {
                return "186.52.170.66"
        }

        hdr := r.Header
        hdrRealIp := hdr.Get("X-Real-Ip")
        hdrForwardedFor := hdr.Get("X-Forwarded-For")

        if hdrRealIp == "" && hdrForwardedFor == "" {
                return ipFromRemote(r.RemoteAddr)
        }

        if hdrForwardedFor != "" {
                // X-Forwarded-For is potentially a list of addresses separated with ","
                parts := strings.Split(hdrForwardedFor, ",")
                for i, p := range parts {
                        parts[i] = strings.TrimSpace(p)
                }
                // TODO: should return first non-local address
                return parts[0]
        }
        return hdrRealIp
}

func geolocate(req *http.Request) geoip.Geolocation {
        var user_addr string
        user_addr = ipAddress(req)
        log.Println(user_addr)

        return <-GetLocation(user_addr)
}

func homepage(w http.ResponseWriter, req *http.Request) {
        var cw *CityWeather
        var coords Coordinates

        location_cookie, err := req.Cookie("mw-location")

        if err == nil {
                log.Println("From Cookie cache")
                parts := strings.Split(location_cookie.Value, "|")
                lat, _ := strconv.ParseFloat(parts[0], 64)
                lng, _ := strconv.ParseFloat(parts[1], 64)
                coords = Coordinates{ lat, lng }
        } else {
                log.Println("From geolocation")
                geo := geolocate(req)
                coords = Coordinates{ geo.Location.Latitude, geo.Location.Longitude }
        }

        city := <-FindByCoords(coords)
        weather := <-GetWeather(city.Coords)

        cw = &CityWeather{ City: city, Weather: weather }

        unit_cookie, err := req.Cookie("mw-unit")
        if err == nil {
                cw.Unit = unit_cookie.Value
        } else {
                cw.Unit = "C"
        }

        if cw.Unit == "F" {
                cw.Weather.Temperature = ((cw.Weather.Temperature*9)/5)+32
        }

        lat, lng := city.Coords.Lat, city.Coords.Lng

        cookie := &http.Cookie{
                Name: "mw-location",
                Value: fmt.Sprintf("%f|%f", lat, lng),
                Path: "/",
        }

        http.SetCookie(w, cookie)

        city_cookie := &http.Cookie{
                Name: "mw-city",
                Value: fmt.Sprintf("%s", city.Name),
                Path: "/",
        }

        http.SetCookie(w, city_cookie)

        t, _ := template.ParseFiles("./website/index.html")
        out, err := json.Marshal(cw)
        cw.JSON = string(out)
        cw.Weather.Temperature = math.Floor(cw.Weather.Temperature)
        err = t.Execute(w, cw)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

type About struct{}

func about(w http.ResponseWriter, req *http.Request) {
        a := &About{}
        t, _ := template.ParseFiles("./website/about.html")
        err := t.Execute(w, a)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}

func isOk(r *http.Request, rm *mux.RouteMatch) bool {
        ref := r.Referer()
        allowed := ref == "http://localhost:12345/" ||
        ref == "http://minimalweather.com/" ||
        ref == "http://www.minimalweather.com/"
        return allowed
}

func Handler() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/weather/{city}", weatherByCity).Methods("GET").MatcherFunc(isOk)
	r.HandleFunc("/weather/{lat}/{lng}", weatherByCoords).Methods("GET").MatcherFunc(isOk)
	r.HandleFunc("/city/{lat}/{lng}", cityByCoords).Methods("GET").MatcherFunc(isOk)

	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./website/")))
        r.HandleFunc("/about", about).Methods("GET")
        r.HandleFunc("/", homepage).Methods("GET")

	return r
}
