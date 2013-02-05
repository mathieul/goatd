app.OverviewCtrl = ($scope, Rpc) ->
  overview = Rpc("Overview", "List")
  overview.List (result) ->
    $scope.resources = result.Rows

app.OverviewCtrl.$inject = ["$scope", "Rpc"]
