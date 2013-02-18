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
        labels:
          title:  "Add a new team"
          action: "Create"

    $scope.editTeam = (team) ->
      $scope.modalTeam.open
        labels:
          title:  "Edit team #{team.name}"
          action: "Update"



])
