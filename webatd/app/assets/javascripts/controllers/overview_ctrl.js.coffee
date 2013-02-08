angular.module("atd").controller("OverviewCtrl", [
  "$scope", "Rpc",
  ($scope, Rpc) ->

    overview = Rpc("Overview", "List")
    overview.List (result) ->
      $scope.resources = result.Rows

])
