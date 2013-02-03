# Controllers

window.TestCtrl = ($scope, $http) ->
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

TestCtrl.$inject = ['$scope', '$http']
