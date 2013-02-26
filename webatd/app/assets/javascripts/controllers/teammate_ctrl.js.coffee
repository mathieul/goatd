angular.module("atd").controller("TeammateCtrl", [
  "$scope", "$q", "ResourceHelper", "Teammate", "Team",
  ($scope, $q, ResourceHelper, Teammate, Team) ->

    deferred = [$q.defer(), $q.defer()]
    teams = Team.index({}, -> deferred[0].resolve())
    teammates = Teammate.index({}, -> deferred[1].resolve())
    bothLoaded = $q.all([deferred[0].promise, deferred[1].promise])
    bothLoaded.then ->
        byUid = {}
        byUid[team.uid] = team.name for team in teams
        teammate.team_name = byUid[teammate.team_uid] for teammate in teammates

    $scope.teammates = teammates

    $scope.addTeammate = ->
      ResourceHelper.openDialog "add-edit-teammate.html", {
        models:
          teams: teams
        labels:
          title:  "Add a new teammate"
          action: "Create"
      } , (result) ->
        if result.action is "save"
          Teammate.create({teammate: result.data}, ResourceHelper.reloader)

    $scope.editTeammate = (teammate) ->
      ResourceHelper.openDialog "add-edit-teammate.html", {
        models:
          teams: teams
          teammate: teammate
        labels:
          title:  "Edit teammate \"#{teammate.name}\""
          action: "Update"
      } , (result) ->
        if result.action is "save"
          Teammate.update({uid: result.data.uid, teammate: result.data}, ResourceHelper.reloader)

    $scope.deleteTeammate = (teammate) ->
      ResourceHelper.confirmDelete "Delete Teammate",
        "Are you sure you want to delete teammate \"#{teammate.name}\"?",
        (choice) ->
          if choice is "delete"
            Teammate.destroy({uid: teammate.uid}, ResourceHelper.reloader)
])
