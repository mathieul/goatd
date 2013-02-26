angular.module("atd").controller("TeamCtrl", [
  "$scope", "ResourceHelper", "Team",
  ($scope, ResourceHelper, Team) ->

    $scope.teams = Team.index()

    $scope.addTeam = ->
      ResourceHelper.openDialog "add-edit-team.html", {
        labels:
          title:  "Add a new team"
          action: "Create"
      } , (result) ->
        if result.action is "save"
          Team.create({team: result.data}, ResourceHelper.reloader)

    $scope.editTeam = (team) ->
      ResourceHelper.openDialog "add-edit-team.html", {
        models:
          team: team
        labels:
          title:  "Edit team \"#{team.name}\""
          action: "Update"
      } , (result) ->
        if result.action is "save"
          Team.update({uid: result.data.uid, team: result.data}, ResourceHelper.reloader)

    $scope.deleteTeam = (team) ->
      ResourceHelper.confirmDelete "Delete Team",
        "Are you sure you want to delete team \"#{team.name}\"?",
        (choice) ->
          if choice is "delete"
            Team.destroy({uid: team.uid}, ResourceHelper.reloader)
])
