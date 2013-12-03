package minimalweather

import (
	"cgl.tideland.biz/asserts"
	"testing"
)

func init() {
	ClearCache()
}

func TestWeatherLookup(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	coords := city.Coordinates{-20, 10}
	weather := <-GetWeather(coords)

	assert.NotNil(weather.Condition, "Should have a condition")
	assert.NotNil(weather.Icon, "Should have an icon")
	assert.NotNil(weather.Temperature, "Should have a temperature")
	assert.NotNil(weather.RainIntensity, "Should have a value")
	assert.NotNil(weather.BringUmbrella, "Should have a value")

	ClearCache()
}
