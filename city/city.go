package city

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/elcuervo/geocoder"
	"github.com/elcuervo/minimalweather/db"
	"github.com/garyburd/redigo/redis"
)

const prefix = "mw:city:"

var c = db.Pool.Get()

type Coordinates geocoder.Coordinates

type City struct {
	Name   string
	Coords Coordinates
}

func ClearCache() {
	pattern := fmt.Sprintf("%s*", prefix)
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
	return fmt.Sprintf("%s%s", prefix, l.Name)
}

func (l *LookupInformation) byCoords() string {
	return fmt.Sprintf("%s%f,%f", prefix, l.Coords.Lat, l.Coords.Lng)
}

func findCity(l LookupInformation, out chan City) {
	cached_city, err := c.Do("GET", l.Key())

	if err != nil {
		panic(err)
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
			city   *geocoder.Location
			coords Coordinates
		)

		switch {
		case l.Name != "":
			city, _ = geocoder.City(l.Name)

		case l.Coords.Lng != 0.0 || l.Coords.Lat != 0.0:
			city, _ = geocoder.Coords(l.Coords.Lat, l.Coords.Lng)

		default:
			panic("OMFG")
		}

		log.Printf("Checking city: %s\n", city.Name)

		coords = Coordinates{city.Coordinates.Lat, city.Coordinates.Lng}
		location := &City{
			Coords: coords,
			Name:   city.Name}

		jsonResponse, _ := json.Marshal(location)

		_, err := c.Do("SET", l.Key(), jsonResponse)

		if err != nil {
			panic(err)
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
