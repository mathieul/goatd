# Modules

angular.module("goatd", ["goatdServices"])
angular.module("goatdServices", ["ng"])

# Application

angular.module("goatd")
  .config(["$routeProvider", ($routeProvider) ->
    $routeProvider
      .when("/",
        templateUrl: "overview.html"
        controller: "OverviewCtrl"
      )
      .when("/teams",
        templateUrl: "teams.html"
        controller: "TeamCtrl"
      )
      .otherwise(redirectTo: "/")
  ])
