angular.module("atd").controller("TeammateCtrl", [
  "$scope", "$route", "Teammate", "BsModal",
  ($scope, $route, Teammate, BsModal) ->

    $scope.teammates = Teammate.index()

    reloader = -> $route.reload()
    $scope.modalTeammate = BsModal "modal-teammate", attributes: ["uid", "name"], save: (attributes) ->
      if attributes.uid?
        Teammate.update(uid: attributes.uid, teammate: attributes, reloader)
      else
        Teammate.create(teammate: attributes, reloader)

    $scope.addTeammate = ->
      $scope.modalTeammate.open
        labels:
          title:  "Add a new teammate"
          action: "Create"

    $scope.editTeammate = (teammate) ->
      $scope.modalTeammate.open
        labels:
          title:  "Edit teammate #{teammate.name}"
          action: "Update"
        values: teammate

    $scope.modalConfirm = BsModal "modal-confirm", attributes: ["uid"], save: (attributes) ->
      Teammate.destroy(uid: attributes.uid, reloader)

    $scope.deleteTeammate = (teammate) ->
      $scope.modalConfirm.open
        labels:
          action: "Delete"
          title: "Delete Teammate"
          question: "you want to delete teammate \"#{teammate.name}\""
        values: teammate
])
