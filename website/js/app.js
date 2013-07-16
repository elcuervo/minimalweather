"use strict";

var weatherAsIcon = function(text) {
  var icon = ")";

  switch(text) {
    case 'wind':                icon = "&#xe020;"; break;
    case 'sleet':               icon = "&#xe010;"; break;
    case 'thunderstorm':        icon = "&#xe00f;"; break;
    case 'snow':
    case 'hail':                icon = "&#xe00c;"; break;
    case 'cloudy':              icon = "&#xe00e;"; break;
    case 'rain':                icon = "&#xe008;"; break;
    case 'fog':                 icon = "&#xe014;"; break;
    case 'clear-day':           icon = "&#xe001;"; break;
    case 'clear-night':         icon = "&#xe002;"; break;
    case 'partly-cloudy-day':   icon = "&#xe000;"; break;
    case 'partly-cloudy-night': icon = "&#xe004;"; break;
    default:                    icon = "&#xe00d;";
  }

  return icon;
};

angular.module("weather", []);
angular.module("weather", []).filter("asWeather", function() { return weatherAsIcon; });

var minimalweather = angular.module("minimalweather", [
  "ngResource", "ngCookies", "ngGeolocation", "weather"
]);

minimalweather.factory("Weather", function($resource, $http) {
  return {
    byName:   $resource("/weather/:city", { city: "@city" }),
    byCoords: $resource("/weather/:lat,:lng", { lat: "@lat", lng: "@lng" })
  }
});

var MainController = function($scope, $resource, $cookieStore, Weather, geolocation) {
  var locateVisitor = function() {
    var currentCity = $cookieStore.get("currentCity");

    if(currentCity && currentCity.coordinates) {
      var lat = currentCity.coordinates.lat;
      var lng = currentCity.coordinates.lng;
      var city = Weather.byCoords.get({ lat: lat, lng: lng });

      console.log("Loaded from cookie cache:", lat, lng);

      $scope.city = city;
    } else {
      geolocation.position().then(function(geo) {
        var lat = geo.coords.latitude;
        var lng = geo.coords.longitude;
        var city = Weather.byCoords.get({ lat: lat, lng: lng });

        console.log("Seek for geolocation:", lat, lng);

        city.$then(function() { $cookieStore.put("currentCity", city) });

        $scope.city = city;
      });
    }
  }

  locateVisitor();

  $scope.clear = function() {
    console.log("Deleted cache");
    $cookieStore.remove("currentCity");
    locateVisitor();
  }

  $scope.search = function() {
    var city = Weather.byName.get({ city: this.city.name });
    console.log("Searching by city name:", this.city.name);

    city.$then(function() { $cookieStore.put("currentCity", city) });

    $scope.city = city;
  }
}

minimalweather.run(function($resource) { });
