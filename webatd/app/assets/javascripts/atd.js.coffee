#
# Modules
#
angular.module("atd", ["atdServices"])
angular.module("atdServices", ["ng"])

#
# Application
#
angular.module("atd").config([
  "$routeProvider",
  ($routeProvider) ->

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