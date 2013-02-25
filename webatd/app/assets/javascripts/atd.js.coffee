#
# Modules
#
angular.module("atd", ["ui.bootstrap", "atdServices"])
angular.module("atdServices", ["ngResource"])

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
      .when("/teammates",
        templateUrl: "teammates.html"
        controller: "TeammateCtrl"
      )
      .otherwise(redirectTo: "/")

])
