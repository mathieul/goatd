angular.module("atd").controller("AddEditResourceCtrl", [
  "$scope", "dialog", "models", "labels",
  ($scope, dialog, models, labels) ->
    $scope.models = models
    $scope.labels = labels

    $scope.action = (action, data) ->
      dialog.close(action: action, data: data)

    $scope.cancel = ->
      dialog.close(action: "cancel")
])
