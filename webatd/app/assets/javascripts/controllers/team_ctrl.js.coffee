angular.module("atd").controller("TeamCtrl", [
  "$scope", "$route", "Team", "BsModal",
  ($scope, $route, Team, BsModal) ->

    $scope.teams = Team.index()

    reloader = -> $route.reload()
    $scope.modalTeam = BsModal "modal-team", attributes: ["uid", "name"], save: (attributes) ->
      if attributes.uid?
        Team.update(uid: attributes.uid, team: attributes, reloader)
      else
        Team.create(team: attributes, reloader)

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
        values: team

    $scope.modalConfirm = BsModal "modal-del-team", attributes: ["uid"], save: (attributes) ->
      Team.destroy(uid: attributes.uid, reloader)

    $scope.deleteTeam = (team) ->
      $scope.modalConfirm.open
        values: team
])
