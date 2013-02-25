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

    # $scope.modalConfirm = BsModal "modal-del-teammate", attributes: ["uid"], save: (attributes) ->
    #   Teammate.destroy(uid: attributes.uid, reloader)

    # $scope.deleteTeammate = (teammate) ->
    #   $scope.modalConfirm.open
    #     values: teammate

    # $scope.openDialog = ->
    #   dialog.open().then (result) ->
    #     console.log("dialog =", result)

    # $scope.openMessageBox = ->
    #   dialog = $dialog.dialog
    #     templateUrl: "add-edit-teammate.html"
    #     controller: "AddEditTeammateCtrl"
    #     modalFade: true
    #     backdropFade: true
    #     resolve:
    #       testing: true
    #       hello: "there"
    #   messageBox = $dialog.messageBox("Delete Teammate",
    #     "Are you sure?",
    #     [
    #       label: "Delete"
    #       cssClass: "btn-primary"
    #       result: "yes"
    #     ,
    #       label: "Cancel"
    #       result: "no"
    #     ]
    #   )
    #   messageBox.open().then (result) -> console.log("messageBox =", result)
])
