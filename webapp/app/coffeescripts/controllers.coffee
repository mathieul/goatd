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


app.OverviewCtrl = ($scope) ->
  $scope.resources = [
    {name: "Team", count: 1}
    {name: "Teammate", count: 5}
    {name: "Queue", count: 3}
    {name: "Task", count: 42}
  ]

  $scope.message = (msg) -> alert(msg)

app.OverviewCtrl.$inject = ['$scope']
