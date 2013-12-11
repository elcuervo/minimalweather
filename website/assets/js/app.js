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
  this.unit = localStorage.getItem("unit");

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

  this.createAppIcon = function(iconFn) {
    var appIcon = this.findOrCreateElement("ios_icon", "apple-touch-icon-precomposed");
    var canvas = document.getElementById("ios_icon_generator");
    var unit, temperature;

    if(this.mw.city.country == "USA" || this.unit == "f" ) {
      temperature = Math.floor(((this.mw.weather.temperature*9)/5)+32);
      unit = "F";
    } else {
      temperature = Math.floor(this.mw.weather.temperature);
      unit = "C";
    }

    canvas.setAttribute('width', 228);
    canvas.setAttribute('height', 228);

    var context = canvas.getContext("2d");
    var gradient = context.createLinearGradient(0, 0, 0, canvas.height);

    gradient.addColorStop(0, '#d55150');
    gradient.addColorStop(1, '#e47d43');

    context.fillStyle = gradient;
    context.fillRect(0, 0, canvas.width, canvas.height);

    iconFn(context);

    if(this.mw.weather.bring_umbrella) icons["umbrella"](context)

    context.fillStyle = "white";
    context.font = "3em Lato"; // temperature
    context.textAlign = "right";

    context.fillText(temperature + "°" + unit, 200, 50);
    context.scale(window.devicePixelRatio, window.devicePixelRatio);

    var data = canvas.toDataURL("image/png");

    appIcon.href = data;
  };
};

MinimalWeather.prototype = {
  generateIcon: function() {
    var icon = this.mw.weather.icon;
    var fn = icons[icon];

    this.createAppIcon(fn);
  }
};

