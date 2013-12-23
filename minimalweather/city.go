package minimalweather

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elcuervo/geocoder"
	"github.com/garyburd/redigo/redis"
)

const city_prefix = "mw:city:"

type Coordinates geocoder.Coordinates

type City struct {
	Name    string      `json:"city"`
	Country string      `json:"country"`
	Coords  Coordinates `json:"coordinates"`
	Error   error       `json:"-"`
}

func ClearCityCache() {
	pattern := fmt.Sprintf("%s*", city_prefix)
	keys, err := redis.Values(c.Do("KEYS", pattern))

	if err != nil {
		log.Println(err)
	}

	for _, key := range keys {
		c.Do("DEL", key)
	}
}

type LookupInformation struct {
	Name   string
	Coords Coordinates

	Primary   string
	Secondary string
}

func (l *LookupInformation) Key() string {
	switch {
	case l.Name != "":
		return l.byName()
	case l.Coords.Lng != 0.0 || l.Coords.Lat != 0.0:
		return l.byCoords()
	}
	return ""
}

func (l *LookupInformation) byName() string {
	return fmt.Sprintf("%s%s", city_prefix, l.Name)
}

func (l *LookupInformation) byCoords() string {
	return fmt.Sprintf("%s%f,%f", city_prefix, l.Coords.Lat, l.Coords.Lng)
}

func findCity(l LookupInformation, out chan City) {
	cached_city, err := c.Do("GET", l.Key())

	if err != nil {
		log.Println(err)
	}

	if cached_city != nil {
		var location City

		str, _ := redis.String(cached_city, nil)
		bytes := []byte(str)
		json.Unmarshal(bytes, &location)

		log.Printf("Loading city: %s\n", location.Name)

		out <- location
	} else {
		var (
			city     *geocoder.Location
			location *City
			err      error
			coords   Coordinates
		)

		if l.Name != "" {
			city, err = geocoder.City(l.Name)
		} else {
			city, err = geocoder.Coords(l.Coords.Lat, l.Coords.Lng)
		}

		if err != nil {
                        log.Println("city.go:94", err)
			location = &City{Name: "Unknown", Error: err}
		} else {
			log.Printf("Checking city: %s\n", city.Name)
			coords = Coordinates{city.Coordinates.Lat, city.Coordinates.Lng}
			location = &City{Coords: coords, Country: city.Country, Name: city.Name}

		}

		json_response, _ := json.Marshal(location)
		_, err = c.Do("SET", l.Key(), json_response)

		if err != nil {
                        log.Println("city.go:106", err)
		}

		out <- *location
	}
}

func FindByCoords(coords Coordinates) chan City {
	city_information := make(chan City)
	lookup := &LookupInformation{Coords: coords}

	go findCity(*lookup, city_information)
	return city_information
}

func FindByName(city_name string) chan City {
	city_information := make(chan City)
	lookup := &LookupInformation{Name: city_name}

	go findCity(*lookup, city_information)
	return city_information
}
