package minimalweather

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/garyburd/redigo/redis"
	forecast "github.com/elcuervo/forecast/v2"
)

const weather_prefix = "mw:weather:"

var api_key = os.Getenv("FORECAST_API_KEY")

type Weather struct {
	Coordinates Coordinates `json:"-"`

	Condition     string  `json:"condition"`
	Icon          string  `json:"icon"`
	Temperature   float64 `json:"temperature"`
	RainIntensity float64 `json:"rain_intensity"`
	BringUmbrella bool    `json:"bring_umbrella"`
}

func ClearWeatherCache() {
	pattern := fmt.Sprintf("%s*", weather_prefix)
	keys, err := redis.Values(c.Do("KEYS", pattern))

	if err != nil {
		log.Println(err)
	}

	for _, key := range keys {
		c.Do("DEL", key)
	}
}

func GetWeather(coords Coordinates) chan Weather {
	city_weather := make(chan Weather)
	key := fmt.Sprintf("%s%f,%f", weather_prefix, coords.Lat, coords.Lng)
	cached_weather, err := c.Do("GET", key)

	if err != nil {
		panic(err)
	}

	go func() {
		if cached_weather != nil {
			var weather Weather

			str, _ := redis.String(cached_weather, nil)
			bytes := []byte(str)
			json.Unmarshal(bytes, &weather)

			log.Printf("Loading weather: %f,%f\n",
				coords.Lat,
				coords.Lng)

			city_weather <- weather
		} else {
			var rain bool

			lat := fmt.Sprintf("%f", coords.Lat)
			lng := fmt.Sprintf("%f", coords.Lng)

			log.Printf("Checking weather: %s,%s\n", lat, lng)

			f, _ := forecast.Get(api_key, lat, lng, "now", forecast.SI)

                        _, err := c.Do("HSET", "mw:stats", "forecast", f.APICalls)
                        if err != nil {
                                log.Println(err)
                        }

			// Look for the next 8 hours. See if it's going to rain
			// at some point
			for _, cond := range f.Hourly.Data[:8] {
				rain = cond.PrecipIntensity > 0.1 || cond.PrecipProbability > 0.6
				if rain {
					break
				}
			}

			weather := &Weather{
				Coordinates:   coords,
				Condition:     f.Currently.Summary,
				Icon:          f.Currently.Icon,
				Temperature:   f.Currently.Temperature,
				RainIntensity: f.Currently.PrecipIntensity,
				BringUmbrella: rain}

			jsonResponse, _ := json.Marshal(weather)
			_, err = c.Do("SETEX", key, 20*60, jsonResponse)

			if err != nil {
				panic(err)
			}

			city_weather <- *weather
		}

	}()

	return city_weather
}
