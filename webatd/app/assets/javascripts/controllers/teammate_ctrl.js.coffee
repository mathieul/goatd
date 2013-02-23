angular.module("atd").controller("TeammateCtrl", [
  "$scope", "$route", "Teammate", "Team", "BsModal",
  ($scope, $route, Teammate, Team, BsModal) ->

    teams = Team.index()
    $scope.teammates = Teammate.index()

    reloader = -> $route.reload()
    $scope.modalTeammate = BsModal "modal-teammate", attributes: ["uid", "name", "team_uid"], save: (attributes) ->
      if attributes.uid?
        Teammate.update(uid: attributes.uid, teammate: attributes, reloader)
      else
        Teammate.create(teammate: attributes, reloader)

    $scope.addTeammate = ->
      $scope.modalTeammate.open
        labels:
          title:  "Add a new teammate"
          action: "Create"
        data:
          teams: teams

    $scope.editTeammate = (teammate) ->
      $scope.modalTeammate.open
        labels:
          title:  "Edit teammate #{teammate.name}"
          action: "Update"
        values: teammate
        data:
          teams: teams

    $scope.modalConfirm = BsModal "modal-del-teammate", attributes: ["uid"], save: (attributes) ->
      Teammate.destroy(uid: attributes.uid, reloader)

    $scope.deleteTeammate = (teammate) ->
      console.log("1- modalConfirm:", $scope.modalConfirm)
      $scope.modalConfirm.open
        values: teammate
      console.log("2- modalConfirm:", $scope.modalConfirm)
])
