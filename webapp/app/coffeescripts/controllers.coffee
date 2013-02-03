# Controllers

app.NavCtrl = ($scope, $location) ->
  $scope.$location = $location

app.NavCtrl.$inject = ["$scope", "$location"]


app.TestCtrl = ($scope) ->
  $scope.items = [
    label: "allo"
    index: 0
  ,
    label: "la"
    index: 1
  ,
    label: "terre"
    index: 2
  ,
    label: "ici"
    index: 3
  ,
    label: "londres"
    index: 4
  ]

  $scope.order = "label"

app.TestCtrl.$inject = ['$scope']


app.JsonCtrl = ($scope, $http) ->
  $http
    .post("/rpc",
      method: "Test.Run"
      params: [{"Name": "Mathieu", "Number": 3}]
      id: "42"
    ,
      headers:
        "Content-Type": "application/json"
    )
    .success (data) ->
      $scope.json = JSON.stringify(data.result)

app.JsonCtrl.$inject = ['$scope', '$http']
