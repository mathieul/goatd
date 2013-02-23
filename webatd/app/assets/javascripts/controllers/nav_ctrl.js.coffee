angular.module("atd").controller("NavCtrl", [
  "$scope", "$location",
  ($scope, $location) ->

    $scope.navs = [
      label: "Overview"
      path: "/"
    ,
      label: "Teams"
      path: "/teams"
    ,
      label: "Teammates"
      path: "/teammates"
    ]

    $scope.activeIfCurrent = (path) ->
      if $location.path() == path then "active" else ""

])
