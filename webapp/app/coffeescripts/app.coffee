# Application
angular
  .module("testing", [])
  .config(["$routeProvider", ($routeProvider) ->
    $routeProvider
      .when('/', {templateUrl: 'home.html'})
      .when('/json', {templateUrl: 'json.html'})
      .otherwise(redirectTo: '/')
  ])

# angular.bootstrap document, ["testing"]
