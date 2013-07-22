module.exports = function(grunt) {
  grunt.loadNpmTasks('grunt-contrib-cssmin');
  grunt.loadNpmTasks("grunt-contrib-uglify");
  grunt.loadNpmTasks('grunt-contrib-watch');

  grunt.initConfig({
    pkg: grunt.file.readJSON("package.json"),

    watch: {
      gruntfile: {
        files: 'Gruntfile.js',
        tasks: ['jshint:gruntfile'],
      },

      src: {
        files: [
          'js/**/*.js',
          'css/**/*.css',

          '!!js/**/*.min.js',
          '!!css/**/*.min.css',
        ],
        tasks: ['build'],
      },
    },

    cssmin: {
      combine: {
        files: {
          "css/<%= pkg.name %>.min.css": [
            "css/pure.css",
            "css/icons.css",
            "css/styles.css"
          ]
        }
      }
    },

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
        dest: "js/<%= pkg.name %>.min.js"
      }
    }
  });

  grunt.registerTask("build", ["uglify", "cssmin"]);
  grunt.registerTask("default", "watch");
}
