module.exports = function(grunt) {
  grunt.loadNpmTasks("grunt-contrib-uglify");

  grunt.initConfig({
    pkg: grunt.file.readJSON("package.json"),

    uglify: {
      options: {
        mangle: false
      },
      build: {
        src: [
          "js/icons.js",
          "js/libs/angular-unstable/angular.js",
          "js/libs/angular-unstable/angular-resource.js",
          "js/libs/angular-unstable/angular-cookies.js",
          "js/libs/angular-geolocation/geolocation.js",
          "js/libs/angular-localstorage/localStorageModule.js",
          "js/konami.js",
          "js/app.js"
        ],
        dest: "js/<%= pkg.name %>.js"
      }
    }
  });

  grunt.registerTask("default", ["uglify"]);
}
