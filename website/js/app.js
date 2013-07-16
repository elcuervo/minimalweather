"use strict";

var weatherAsIcon = function(text) {
  var icon = ")";

  switch(text) {
    case 'loading':             icon = "("; break;

    case 'wind':                icon = "F"; break;
    case 'sleet':               icon = "$"; break;
    case 'thunderstorm':        icon = "&"; break;
    case 'snow':
    case 'hail':                icon = "#"; break;
    case 'cloudy':              icon = "%"; break;
    case 'rain':                icon = "8"; break;
    case 'fog':                 icon = "M"; break;
    case 'clear-day':           icon = "1"; break;
    case 'clear-night':         icon = "2"; break;
    case 'partly-cloudy-day':   icon = "3"; break;
    case 'partly-cloudy-night': icon = "4"; break;
    default:                    icon = "1";
  }

  return icon;
};

angular.module("weather", []);
angular.module("weather", []).filter("asWeather", function() { return weatherAsIcon; });

var minimalweather = angular.module("minimalweather", [
  "ngResource", "ngCookies", "ngGeolocation", "weather"
]);

minimalweather.factory("Weather", function($resource) {
  return {
    byName:   $resource("/weather/:city", { city: "@city" }),
    byCoords: $resource("/weather/:lat/:lng", { lat: "@lat", lng: "@lng" })
  }
})

var MainController = function($scope, $resource, $cookieStore, Weather, geolocation) {
  var locateVisitor = function() {
    var coords = $cookieStore.get("coordinates");

    if(coords) {
      var lat = coords.lat;
      var lng = coords.lng;

      console.log("Loaded from cookie cache:", lat, lng);

      $scope.city = Weather.byCoords.get({ lat: lat, lng: lng });
      $scope.citySearch = $scope.city.name;
    } else {
      geolocation.position().then(function(geo) {
        var lat = geo.coords.latitude;
        var lng = geo.coords.longitude;

        console.log("Seek for geolocation:", lat, lng);

        $cookieStore.put("coordinates", { lat: lat, lng: lng });

        $scope.city = Weather.byCoords.get({ lat: lat, lng: lng });
        $scope.citySearch = $scope.city.name;
      });
    }
  }

  $scope.city = { weather: { icon: "loading" } };

  locateVisitor();

  $scope.clear = function() {
    console.log("Deleted cache");
    $cookieStore.remove("coordinates");
    locateVisitor()
  }

  $scope.search = function() {
    console.log("Searching by city name:", this.citySearch);
    var city = Weather.byName.get({ city: this.citySearch });
    $cookieStore.put("coordinates2", "test");

    $scope.city = city;
  }
}

minimalweather.run(function($resource) { });
