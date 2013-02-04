# Controllers

app.NavCtrl = ($scope, $location) ->
  $scope.navs = [
    label: "Overview"
    path: "/"
  ,
    label: "TODO"
    path: "/todo"
  ]

  $scope.activeIfCurrent = (path) ->
    if $location.path() == path then "active" else ""

app.NavCtrl.$inject = ["$scope", "$location"]


app.OverviewCtrl = ($scope, Rpc) ->
  overview = Rpc("Overview", "List")
  overview.List (result) ->
    $scope.resources = result.Rows
  $scope.message = (msg) -> alert(msg)

app.OverviewCtrl.$inject = ["$scope", "Rpc"]
