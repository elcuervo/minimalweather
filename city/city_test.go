package city

import (
	"cgl.tideland.biz/asserts"
	"testing"
)

func init() {
	ClearCache()
}

func TestFindByName(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	coords := Coordinates{-34.8836111, -56.1819444}
	city := <-FindByName("Montevideo")

	assert.Equal(city.Name, "Montevideo", "Wrong city name")
	assert.Equal(city.Coords, coords, "Wrong city coordinates")

	ClearCache()
}

func TestFindByCoords(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	coords := Coordinates{10.8112436, 20.5215459}
	city := <-FindByCoords(coords)

	assert.Equal(city.Name, "Barh Azoum", "Wrong city name")
	assert.Equal(city.Coords, coords, "Wrong city coordinates")

	ClearCache()
}

func TestRetreiveFromCache(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	coords := Coordinates{-34.8836111, -56.1819444}
	city := <-FindByName("Montevideo")

	assert.Equal(city.Name, "Montevideo", "Wrong city name")
	assert.Equal(city.Coords, coords, "Wrong city coordinates")

	city = <-FindByName("Montevideo")
	assert.Equal(city.Coords, coords, "Wrong city coordinates")

	ClearCache()
}

func TestWrongCity(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	city := <-FindByName("Blarjs223")
	assert.NotNil(city.Error, "The city should not be right")
}
