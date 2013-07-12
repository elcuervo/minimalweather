angular.module('ngGeolocation',[])
  .constant('options',{})
  .factory('geolocation',
        ["$q","$rootScope","options",
function ($q , $rootScope , options){
  return {
    position: function () {
      var deferred = $q.defer()
      navigator.geolocation.getCurrentPosition(function (pos) {
        $rootScope.$apply(function () {
          deferred.resolve(angular.copy(pos))
        })
      }, function (error) {
        $rootScope.$apply(function () {
          deferred.reject(error)
        })
      },options)
      return deferred.promise
    }
  }
}])