package main

import (
        "net/http"
        "log"
        "os"
        "fmt"
        "strconv"
        "encoding/json"
        "github.com/elcuervo/geocoder"
        "github.com/bmizerany/pat"
        forecast "github.com/mlbright/forecast/v2"
)

type Weather struct {
        City string `json:"city"`
        Coordinates geocoder.Coordinates `json:"coordinates"`

        Condition       string  `json:"condition"`
        Temperature     float64 `json:"temperature"`
        RainProbability float64 `json:"rain_probability"`
}

func FindCityCoords(w http.ResponseWriter, req *http.Request) {
        city_name := req.URL.Query().Get(":city")
        city, _ := geocoder.City(city_name)

        res, _ := json.Marshal(cityWeather(city))
        w.Write(res)
}

func WeatherHandler(w http.ResponseWriter, req *http.Request) {
        lat, _ := strconv.ParseFloat(req.URL.Query().Get(":lat"), 64)
        lng, _ := strconv.ParseFloat(req.URL.Query().Get(":lng"), 64)
        city, _ := geocoder.Coords(lat, lng)

        log.Println(cityWeather(city))
}

func cityWeather(city *geocoder.Location) *Weather {
        lat := fmt.Sprintf("%f", city.Coordinates.Lat)
        long := fmt.Sprintf("%f", city.Coordinates.Lng)
        key := os.Getenv("FORECAST_API_KEY")

        f := forecast.Get(key, lat, long, "now")
        return &Weather{
                City: city.Name,
                Coordinates: city.Coordinates,
                Condition: f.Currently.Summary,
                Temperature: f.Currently.Temperature,
                RainProbability: f.Currently.PrecipProbability * 100}
}

func main() {
        m := pat.New()
        port := ":12345"

        m.Get("/weather/:lat,:lng", http.HandlerFunc(WeatherHandler))
        m.Get("/weather/:city", http.HandlerFunc(FindCityCoords))
        http.Handle("/", m)

        log.Println("Listening in", port)
        err := http.ListenAndServe(port, nil)

        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}
