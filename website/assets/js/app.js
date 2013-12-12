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

var MinimalWeather = function(json) {
  this.mw = JSON.parse(json);

  this.refreshLocationTo = function(lat, lng) {
    document.getElementById("city_name").textContent = "Relocating";
    document.getElementById("condition_icon").className = "icon-relocating";

    this.cookieMonster.set("mw-location", lat + "|" + lng);
    location.reload();
  };

  this.findOrCreateElement = function(id, rel) {
    var iosIcon = document.getElementById(id);

    if(!!iosIcon) {
      return iosIcon;
    } else {
      var link = document.createElement("link");

      link.id = id;
      link.rel = rel;

      document.head.appendChild(link);

      return link;
    }
  }

  this.cookieMonster = {
    get: function(key) {
      var data = this.all();
      return data[key];
    },

    all: function() {
      var cookies = document.cookie.split(';');
      var data = {};

      for(var i = 0; i < cookies.length; i++) {
        var keysAndValues = cookies[i].split('=');
        data[keysAndValues[0].replace( /^\s+|\s+$/g, '' )] = keysAndValues[1];
      }
      return data;
    },

    del: function(key) {
      this.set(key, '', '; expires=Thu, 01 Jan 1970 00:00:01 GMT;' )
    },

    set: function(key, value, expires) {
      if(!expires) {
        var date = new Date();
        date.setTime(date.getTime() + (60*24*60*60*1000));
        var expires = "; expires=" + date.toGMTString();
      }

      document.cookie = key + '=' + value + expires +  'path=/'
    }
  };

  this.usesFarenheit = function() {
    var cookieCache = this.cookieMonster.get("mw-unit").toUpperCase() == "F";
    return cookieCache || this.mw.city.country == "USA";
  };

  this.createAppIcon = function(iconFn) {
    var appIcon = this.findOrCreateElement("ios_icon", "apple-touch-icon-precomposed");
    var canvas = document.getElementById("ios_icon_generator");
    var temperature = Math.floor(this.mw.weather.temperature);
    var unit = this.mw.unit;

    canvas.setAttribute('width', 228);
    canvas.setAttribute('height', 228);

    var context = canvas.getContext("2d");
    var gradient = context.createLinearGradient(0, 0, 0, canvas.height);

    if(this.mw.cold) {
      gradient.addColorStop(0, '#1e5799');
      gradient.addColorStop(1, '#7db9e8');
    } else {
      gradient.addColorStop(0, '#d55150');
      gradient.addColorStop(1, '#e47d43');
    }

    context.fillStyle = gradient;
    context.fillRect(0, 0, canvas.width, canvas.height);

    iconFn(context);

    if(this.mw.weather.bring_umbrella) icons["umbrella"](context)

    context.fillStyle = "white";
    context.font = "3em Helvetica"; // temperature
    context.textAlign = "right";

    context.fillText(temperature + "°" + unit, 200, 50);
    context.scale(window.devicePixelRatio, window.devicePixelRatio);

    var data = canvas.toDataURL("image/png");

    appIcon.href = data;
  };

  var self = this;

  new Konami(function() {
    self.refreshLocationTo("-27.1167", "-109.3667")
    self.cookieMonster.set("mw-easter", true)
  });
};

MinimalWeather.prototype = {
  generateIcon: function() {
    var icon = this.mw.weather.icon;
    var fn = icons[icon];

    this.createAppIcon(fn);
  },

  bindUnits: function(units) {
    var self = this;
    var changeUnit = function() {
      self.cookieMonster.set("mw-unit", this.id.toUpperCase());
    }
    for(i in units) {
      if(units[i]) units[i].addEventListener("click", changeUnit);
    }
  },

  bindRefresh: function(button) {
    var self = this;
    button.addEventListener("click", function() {
      self.cookieMonster.del("mw-city");
      self.cookieMonster.del("mw-location");
    })
  },

  geolocate: function() {
    var self = this;
    navigator.geolocation.getCurrentPosition(function(position) {
      var cookieContent = self.cookieMonster.get("mw-city")
      var cookieCity = decodeURIComponent(cookieContent.replace(/\+/g, ' '));

      var lat = position.coords.latitude;
      var lng = position.coords.longitude;

      $.getJSON('/city/' + lat + '/' + lng, function(response){
        if(response.city == undefined) return;


        if(self.cookieMonster.get("mw-easter")) {
          setTimeout(function() {
            self.cookieMonster.del("mw-easter")
            self.refreshLocationTo(lat, lng);
          }, 5000);
        }

        if(response.city != cookieCity) {
          console.log("Ok, you moved from " + cookieCity +  " to " + response.city + ". Let's find you current weather");
          self.refreshLocationTo(lat, lng)
        }

      })
    });
  }
};

