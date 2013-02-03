# Application

window.app = {}

angular
  .module("app.testing", [])
  .config(["$routeProvider", ($routeProvider) ->
    $routeProvider
      .when('/', {templateUrl: 'home.html', controller: 'app.TestCtrl'})
      .when('/json', {templateUrl: 'json.html', controller: 'app.JsonCtrl'})
      .otherwise(redirectTo: '/')
  ])
