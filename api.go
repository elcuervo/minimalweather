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
        Icon            string  `json:"icon"`
        Temperature     float64 `json:"temperature"`
        RainIntensity   float64 `json:"rain_intensity"`
        BringUmbrella   bool    `json:"bring_umbrella"`
}

var (
        api_key = os.Getenv("FORECAST_API_KEY")
        city_weather = make(chan *Weather)
        city_information = make(chan *geocoder.Location)
)

func weatherByCity(w http.ResponseWriter, req *http.Request) {
        city_name := req.URL.Query().Get(":city")

        go func() {
                log.Println("Checking city")
                city, _ := geocoder.City(city_name)
                city_information <- city
        }()

        city := <- city_information
        go cityWeather(city.Coordinates)

        weather := <-city_weather
        weather.City = city.Name
        jsonResponse, _ := json.Marshal(weather)

        w.Write(jsonResponse)
}

func weatherByCoords(w http.ResponseWriter, req *http.Request) {
        lat, _ := strconv.ParseFloat(req.URL.Query().Get(":lat"), 64)
        lng, _ := strconv.ParseFloat(req.URL.Query().Get(":lng"), 64)
        coords := &geocoder.Coordinates{lat, lng}

        go func() {
                log.Println("Checking city")
                city, _ := geocoder.Coords(lat, lng)
                city_information <- city
        }()

        go cityWeather(*coords)

        city := <- city_information
        weather := <-city_weather

        weather.City = city.Name
        jsonResponse, _ := json.Marshal(weather)

        w.Write(jsonResponse)
}

func cityWeather(coords geocoder.Coordinates) {
        log.Println("Checking weather")
        lat := fmt.Sprintf("%f", coords.Lat)
        lng := fmt.Sprintf("%f", coords.Lng)

        f := forecast.Get(api_key, lat, lng, "now")
        future_condition := f.Hourly.Data[6]
        raining := f.Currently.PrecipIntensity > 0.1   ||
                   f.Currently.PrecipProbability > 0.6 ||
                   future_condition.PrecipProbability > 0.6

        city_weather <- &Weather{
                Coordinates: coords,
                Condition: f.Currently.Summary,
                Icon: f.Currently.Icon,
                Temperature: f.Currently.Temperature,
                RainIntensity: f.Currently.PrecipIntensity,
                BringUmbrella: raining}
}

func main() {
        m := pat.New()
        port := ":12345"

        m.Get("/weather/:city", http.HandlerFunc(weatherByCity))
        m.Get("/weather/:lat,:lng", http.HandlerFunc(weatherByCoords))
        http.Handle("/", m)

        log.Println("Listening in", port)
        err := http.ListenAndServe(port, nil)

        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}
