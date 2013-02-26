angular.module("atdServices").factory("ResourceHelper", [
  "$dialog", "$route",
  ($dialog, $route) ->
    {
      openDialog: (template, options, done) ->
        options.models ||= {}
        options.labels ||= {}

        dialog = $dialog.dialog
          templateUrl: template
          controller: "AddEditResourceCtrl"
          modalFade: true
          backdropFade: true
          resolve: options
        dialog.open().then(done)

      confirmDelete: (title, question, done) ->
        messageBox = $dialog.messageBox title, question,
          [
            label: "Delete"
            cssClass: "btn-primary"
            result: "delete"
          ,
            label: "Cancel"
            result: "cancel"
          ]
        messageBox.open().then(done)

      reloader: -> $route.reload()
    }
])
