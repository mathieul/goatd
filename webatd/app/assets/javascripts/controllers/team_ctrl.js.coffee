angular.module("atd").controller("TeamCtrl", [
  "$scope", "$window", "Team", "BsModal",
  ($scope, $window, Team, BsModal) ->

    $scope.message = (msg) -> $window.alert(msg)
    $scope.teams = Team.query()

    modal = BsModal("modal-team")
    $scope.modalTeam =
      submit:
        -> modal.submit()
      save: ->
        console.log("modalTeam.save(): TODO", arguments)
        modal.close()
        false
    $scope.addTeam = ->
      $scope.modalTeam.title = "Add a new team"
      $scope.modalTeam.actionLabel = "Create"
      modal.open()

])
