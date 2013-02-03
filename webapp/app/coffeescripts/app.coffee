# Application

window.app = {}

angular.module("app.goatd", [])
  .config(["$routeProvider", ($routeProvider) ->
    $routeProvider
      .when('/', {templateUrl: 'overview.html', controller: 'app.OverviewCtrl'})
      .otherwise(redirectTo: '/')
  ])
