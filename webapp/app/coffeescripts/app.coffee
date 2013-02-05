# Application

window.app = {}

angular.module("app.goatd", ["app.goatdServices"])
  .config(["$routeProvider", ($routeProvider) ->
    $routeProvider
      .when("/",
        templateUrl: "overview.html"
        controller: "app.OverviewCtrl"
      )
      .when("/teams",
        templateUrl: "teams.html"
        controller: "app.TeamCtrl"
      )
      .otherwise(redirectTo: "/")
  ])
