# Application

window.app = {}

angular.module("app.goatd", ["app.goatdServices"])
  .config(["$routeProvider", ($routeProvider) ->
    $routeProvider
      .when("/", {templateUrl: "overview.html", controller: "app.OverviewCtrl"})
      .otherwise(redirectTo: "/")
  ])
