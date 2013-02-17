angular.module("atd").controller("TeamCtrl", [
  "$scope", "$route", "Team", "BsModal",
  ($scope, $route, Team, BsModal) ->

    $scope.teams = Team.query()

    $scope.modalTeam = BsModal "modal-team", save: (attributes) ->
      Team.save(attributes, ->
        console.log("saved:", attributes)
        $route.reload()
      )

    $scope.addTeam = ->
      $scope.modalTeam.open
        title:  "Add a new team"
        action: "Create"

])
