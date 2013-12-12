package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/elcuervo/minimalweather/pages"
	mw "github.com/elcuervo/minimalweather/minimalweather"
	"log"
	"net/http"
        "html/template"
	"strconv"
)

func outputWeatherAsJSON(current_city mw.City, current_weather mw.Weather) []byte {
	city_weather := &pages.CityWeather{
		City:        current_city,
		Weather:     current_weather}

	out, _ := json.Marshal(city_weather)

	return out
}

func weatherByCity(w http.ResponseWriter, req *http.Request) {
	city_name := mux.Vars(req)["city"]

	log.Println("By Name:", city_name)

	current_city := <-mw.FindByName(city_name)
	current_weather := <-mw.GetWeather(current_city.Coords)

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

	coords := mw.Coordinates{ lat, lng }
	current_city := mw.FindByCoords(coords)
	current_weather := <-mw.GetWeather(coords)

	out := outputWeatherAsJSON(<-current_city, current_weather)

	w.Write(out)
}

func cityByCoords(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	lat, _ := strconv.ParseFloat(vars["lat"], 64)
	lng, _ := strconv.ParseFloat(vars["lng"], 64)

	log.Println("city By Coords:", lat, lng)

	coords := mw.Coordinates{ lat, lng }
	current_city := <-mw.FindByCoords(coords)

        out, _ := json.Marshal(current_city)
	w.Write(out)
}

func homepage(w http.ResponseWriter, req *http.Request) {
        home := pages.NewHomepage(w, req)
        home.Render()
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
        var ref = r.Referer()

        allowed := ref == "http://localhost:12345/" ||
        ref == "http://minimalweather.com/" ||
        ref == "http://nimbus.minimalweather.com/" ||
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
