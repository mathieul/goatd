# Controllers

app.NavCtrl = ($scope, $location) ->
  $scope.navs = [
    label: "Home"
    path: "/"
  ,
    label: "TODO"
    path: "/todo"
  ]

  $scope.activeIfCurrent = (path) ->
    if $location.path() == path then "active" else ""

app.NavCtrl.$inject = ["$scope", "$location"]


app.HomeCtrl = ($scope) ->

app.HomeCtrl.$inject = ['$scope']
