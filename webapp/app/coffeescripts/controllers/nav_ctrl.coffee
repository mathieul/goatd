app.NavCtrl = ($scope, $location) ->
  $scope.navs = [
    label: "Overview"
    path: "/"
  ,
    label: "Teams"
    path: "/teams"
  ]

  $scope.activeIfCurrent = (path) ->
    if $location.path() == path then "active" else ""

app.NavCtrl.$inject = ["$scope", "$location"]
