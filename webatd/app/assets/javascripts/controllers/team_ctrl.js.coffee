angular.module("atd").controller("TeamCtrl", [
  "$scope", "$window", "Team",
  ($scope, $window, Team) ->

    $scope.message = (msg) -> $window.alert(msg)
    $scope.teams = Team.query()

])
