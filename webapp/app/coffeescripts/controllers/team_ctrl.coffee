app.TeamCtrl = ($scope, Rpc) ->
  team = Rpc("Team", "List")
  team.List (result) ->
    $scope.teams = result.Rows
    $scope.noTeams = $scope.teams.length is 0

app.TeamCtrl.$inject = ["$scope", "Rpc"]
