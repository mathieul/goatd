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


app.OverviewCtrl = ($scope, $http) ->
  $http
    .post("rpc",
      method: "Overview.List"
      params: []
      id:     (new Date).getTime()
    ,
      headers:
        "Content-Type": "application/json"
    )
    .success (data) ->
      $scope.resources = data.result.Rows

  $scope.message = (msg) -> alert(msg)

app.OverviewCtrl.$inject = ['$scope', '$http']
