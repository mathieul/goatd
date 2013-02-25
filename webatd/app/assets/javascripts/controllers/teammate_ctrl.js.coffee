openDialog = ($dialog, options, done) ->
  options.models ||= {}
  options.labels ||= {}

  dialog = $dialog.dialog
    templateUrl: "add-edit-teammate.html"
    controller: "AddEditResourceCtrl"
    modalFade: true
    backdropFade: true
    resolve: options
  dialog.open().then(done)

angular.module("atd").controller("TeammateCtrl", [
  "$scope", "$route", "$dialog", "Teammate", "Team",
  ($scope, $route, $dialog, Teammate, Team) ->

    teams = Team.index()
    $scope.teammates = Teammate.index()

    reloader = -> $route.reload()

    $scope.addTeammate = ->
      openDialog $dialog, {
        models:
          teams: teams
        labels:
          title:  "Add a new teammate"
          action: "Create"
      } , (result) ->
        if result.action is "save"
          Teammate.create(teammate: result.data, reloader)

    $scope.editTeammate = (teammate) ->
      openDialog $dialog, {
        models:
          teams: teams
          teammate: teammate
        labels:
          title:  "Edit teammate \"#{teammate.name}\""
          action: "Update"
      } , (result) ->
        if result.action is "save"
          Teammate.update(uid: result.data.uid, teammate: result.data, reloader)

    $scope.deleteTeammate = (teammate) ->
      messageBox = $dialog.messageBox "Delete Teammate",
        "Are you sure you want to delete teammate \"#{teammate.name}\"?",
        [
          label: "Delete"
          cssClass: "btn-primary"
          result: "delete"
        ,
          label: "Cancel"
          result: "cancel"
        ]
      messageBox.open().then (choice) ->
        if choice is "delete"
          Teammate.destroy(uid: teammate.uid, reloader)
])
