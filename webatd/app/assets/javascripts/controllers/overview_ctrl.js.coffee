angular.module("atd").controller("OverviewCtrl", [
  "$scope", "Rpc",
  ($scope, Rpc) ->

    $scope.resources = []
    overview = Rpc("Overview", "List")
    overview.List (result) ->
      $scope.resources = result.Rows

])
