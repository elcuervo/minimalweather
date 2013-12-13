package pages

import (
	"encoding/json"
	"fmt"
	"github.com/elcuervo/geoip"
	mw "github.com/elcuervo/minimalweather/minimalweather"
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Homepage struct {
	w  http.ResponseWriter
	r  *http.Request
	cw *CityWeather
}

func (h *Homepage) getCoords() mw.Coordinates {
	var coords mw.Coordinates
	location_cookie, err := h.r.Cookie("mw-location")

	if err == nil {
		log.Println("From Cookie cache")
		parts := strings.Split(location_cookie.Value, "|")
		lat, _ := strconv.ParseFloat(parts[0], 64)
		lng, _ := strconv.ParseFloat(parts[1], 64)
		coords = mw.Coordinates{lat, lng}
	} else {
		log.Println("From geolocation")
		geo := h.geolocate()
		coords = mw.Coordinates{geo.Location.Latitude, geo.Location.Longitude}
	}

	return coords
}

func (h *Homepage) ipFromRemote() string {
	index := strings.LastIndex(h.r.RemoteAddr, ":")
	if index == -1 {
		return h.r.RemoteAddr
	}
	return h.r.RemoteAddr[:index]
}

func (h *Homepage) ipAddress() string {
	development := os.Getenv("DEVELOPMENT")
	if development != "" {
		return "186.52.170.66"
	}

	hdr := h.r.Header
	hdrRealIp := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")

	if hdrRealIp == "" && hdrForwardedFor == "" {
		return h.ipFromRemote()
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

func (h *Homepage) geolocate() geoip.Geolocation {
	var user_addr string

	user_addr = h.ipAddress()
	log.Println(user_addr)

	return <-mw.GetLocation(user_addr)
}

func (h *Homepage) handleUnit() {
	unit_cookie, err := h.r.Cookie("mw-unit")
	if err == nil {
		h.cw.Unit = unit_cookie.Value
	} else {
                if h.cw.City.Country == "US" {
                        h.cw.Unit = "F"
                } else {
                        h.cw.Unit = "C"
                }
	}

	// Based on @chadot knowladge < 17 is minimun to confort temperature
	if h.cw.Weather.Temperature < 17 {
		h.cw.Cold = true
	} else {
		h.cw.Cold = false
	}

	if h.cw.Unit == "F" {
		h.cw.Weather.Temperature = ((h.cw.Weather.Temperature * 9) / 5) + 32
	}

	h.cw.Celsius = h.cw.Unit == "C"
}

func (h *Homepage) saveCityCache(city mw.City) {
	cookie := &http.Cookie{
		Name:  "mw-location",
		Value: fmt.Sprintf("%f|%f", city.Coords.Lat, city.Coords.Lng),
		Path:  "/",
	}

	http.SetCookie(h.w, cookie)

	city_cookie := &http.Cookie{
		Name:  "mw-city",
		Value: fmt.Sprintf("%s", url.QueryEscape(city.Name)),
		Path:  "/",
	}

	http.SetCookie(h.w, city_cookie)
}

func (h *Homepage) Render() {
	coords := h.getCoords()
	city := <-mw.FindByCoords(coords)
	weather := <-mw.GetWeather(city.Coords)

	h.cw = &CityWeather{City: city, Weather: weather}

	h.handleUnit()
	h.saveCityCache(city)

	t, _ := template.ParseFiles("./website/index.html")
	out, err := json.Marshal(h.cw)
	h.cw.JSON = string(out)
	h.cw.Weather.Temperature = math.Floor(h.cw.Weather.Temperature)
	err = t.Execute(h.w, h.cw)

	if err != nil {
		http.Error(h.w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHomepage(w http.ResponseWriter, req *http.Request) *Homepage {
	home := new(Homepage)
	home.w = w
	home.r = req

	return home
}
