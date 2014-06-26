package routes

import (
	"net/http"
	"time"

	"github.com/elcuervo/minimalweather/pages"
	"github.com/gorilla/mux"
)

var api = new(pages.API)
var about = new(pages.About)

func homepage(w http.ResponseWriter, req *http.Request) {
	home := pages.NewHomepage(w, req)
	home.Render()
}

func isOk(r *http.Request, rm *mux.RouteMatch) bool {
	var ref = r.Referer()

	allowed := ref == "http://localhost:12345/" ||
		ref == "http://nevermore.local:12345/" ||
		ref == "http://minimalweather.com/" ||
		ref == "http://nimbus.minimalweather.com/" ||
		ref == "http://www.minimalweather.com/"

	return allowed
}

func app(w http.ResponseWriter, req *http.Request) {
	expire := time.Now().Add(60 * 24 * 60 * 60 * 1000)
	cookie := http.Cookie{Name: "mw-desktop", Value: "1", Path: "/", Expires: expire, MaxAge: 86400}

	http.SetCookie(w, &cookie)
	http.Redirect(w, req, "/", 301)
}

func Handler() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/weather/{city}", api.WeatherByCity).Methods("GET").MatcherFunc(isOk)
	r.HandleFunc("/weather/{lat}/{lng}", api.WeatherByCoords).Methods("GET").MatcherFunc(isOk)
	r.HandleFunc("/city/{lat}/{lng}", api.CityByCoords).Methods("GET").MatcherFunc(isOk)

	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./website/")))
	r.HandleFunc("/about", about.Render).Methods("GET")
	r.HandleFunc("/app", app).Methods("GET")
	r.HandleFunc("/", homepage).Methods("GET")

	return r
}
