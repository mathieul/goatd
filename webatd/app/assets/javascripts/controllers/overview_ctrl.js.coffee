angular.module("atd").controller("OverviewCtrl", [
  "$scope", "$http",
  ($scope, $http) ->

    $http
      .get("overview", {}, headers: "Content-Type": "application/json")
      .success (result) ->
        $scope.resources = result.rows

])
