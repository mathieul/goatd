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

    $scope.modalConfirm = BsModal "modal-confirm", attributes: ["uid"], save: (attributes) ->
      console.log(attributes)

    $scope.deleteTeam = (team) ->
      $scope.modalConfirm.open
        labels:
          action: "Delete"
          title: "Delete Team"
          question: "you want to delete team \"#{team.name}\""
        values: team
])
