# Controllers
window.TestCtrl = ($scope, $http) ->
  $scope.items = [
    label: "allo"
    index: 2
  ,
    label: "la"
    index: 1
  ,
    label: "terre"
    index: 0
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