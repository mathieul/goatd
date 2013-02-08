angular.module("atd").controller("TeamCtrl", [
  "$scope", "Rpc",
  ($scope, Rpc) ->

    $scope.teams = []
    team = Rpc("Team", "List")
    team.List (result) ->
      $scope.teams = result.Rows

])
