angular.module("atd").controller("TeamCtrl", [
  "$scope", "$route", "Team", "BsModal",
  ($scope, $route, Team, BsModal) ->

    $scope.teams = Team.index()

    $scope.modalTeam = BsModal "modal-team", attributes: ["name"], save: (attributes) ->
      Team.create(attributes, ->
        $route.reload()
      )

    $scope.addTeam = ->
      $scope.modalTeam.open
        title:  "Add a new team"
        action: "Create"

])
