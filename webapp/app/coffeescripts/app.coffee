# Application

window.app = {}

angular
  .module("app.goatd", [])
  .config(["$routeProvider", ($routeProvider) ->
    $routeProvider
      .when('/', {templateUrl: 'home.html', controller: 'app.HomeCtrl'})
      .otherwise(redirectTo: '/')
  ])
