"use strict";

angular.module("weather", []);

angular.module("weather", []).filter("asWeather", function() {
  return function(text) {
    var nightTime = false;
    var icon = ")";

    //, wind, , ,
    //, , or tornado

    switch(text) {
      case 'sleet':
        icon = "X";
        break;
      case 'thunderstorm':
        icon = "Z";
        break;
      case 'snow':
      case 'hail':
        icon = "W";
        break;
      case 'cloudy':
        icon = "Y";
        break;
      case 'rain':
        icon = "R";
        break;
      case 'fog':
        icon = "E";
        break;
      case 'clear-day':
        icon = "B";
        break;
      case 'clear-night':
        icon = "C";
        break;
      case 'partly-cloudy-day':
        icon = "H";
        break;
      case 'partly-cloudy-night':
        icon = "I";
        break;
      default:
        icon = "B";

    }
    return icon;
  };
});


var minimalweather = angular.module("minimalweather", [
    "ngResource", "ngCookies", "ngGeolocation", "weather"
]);

minimalweather.factory("Weather", function($resource) {
  return {
    byName: $resource("/weather/:city", { city: "@city" }),
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
