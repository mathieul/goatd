openDialog = ($dialog, options, done) ->
  options.models ||= {}
  options.labels ||= {}

  dialog = $dialog.dialog
    templateUrl: "add-edit-team.html"
    controller: "AddEditResourceCtrl"
    modalFade: true
    backdropFade: true
    resolve: options
  dialog.open().then(done)

angular.module("atd").controller("TeamCtrl", [
  "$scope", "$route", "$dialog", "Team",
  ($scope, $route, $dialog, Team) ->

    $scope.teams = Team.index()

    reloader = -> $route.reload()
    $scope.addTeam = ->
      openDialog $dialog, {
        labels:
          title:  "Add a new team"
          action: "Create"
      } , (result) ->
        if result.action is "save"
          Team.create(team: result.data, reloader)

    $scope.editTeam = (team) ->
      openDialog $dialog, {
        models:
          team: team
        labels:
          title:  "Edit team \"#{team.name}\""
          action: "Update"
      } , (result) ->
        if result.action is "save"
          Team.update(uid: result.data.uid, team: result.data, reloader)

    # $scope.modalConfirm = BsModal "modal-del-team", attributes: ["uid"], save: (attributes) ->
    #   Team.destroy(uid: attributes.uid, reloader)

    # $scope.deleteTeam = (team) ->
    #   $scope.modalConfirm.open
    #     values: team
])
