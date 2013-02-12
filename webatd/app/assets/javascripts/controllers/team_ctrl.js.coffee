angular.module("atd").controller("TeamCtrl", [
  "$scope", "$http", "$window",
  ($scope, $http, $window) ->

    $scope.message = (msg) ->
      $window.alert(msg)

    $http
      .get("teams", {}, headers: "Content-Type": "application/json")
      .success (result) ->
        $scope.teams = result.teams
])
