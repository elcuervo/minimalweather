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
          'assets/js/**/*.js',
          'assets/css/**/*.css',

          '!!assets/js/**/*.min.js',
          '!!assets/css/**/*.min.css',
        ],
        tasks: ['build'],
      },
    },

    cssmin: {
      combine: {
        files: {
          "assets/css/<%= pkg.name %>.min.css": [
            "assets/css/pure.css",
            "assets/css/icons.css",
            "assets/css/styles.css"
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
          "assets/js/icons.js",
          "assets/js/konami.js",
          "assets/js/app.js"
        ],
        dest: "assets/js/<%= pkg.name %>.min.js"
      }
    }
  });

  grunt.registerTask("build", ["uglify", "cssmin"]);
  grunt.registerTask("default", "watch");
}
