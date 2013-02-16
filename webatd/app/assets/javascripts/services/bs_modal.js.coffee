class ModalManager
  constructor: (id) -> @sel = "##{id}"

  open: (callback) ->
    $(@sel)
      .modal("show")
      .one "shown", (event) ->
        $(event.target).find("form input[type=text]:visible:first")[0].focus()

  close: ->
    $(@sel).modal("hide")

  submit: ->
    $("#{@sel} form:first").submit()

angular.module("atdServices").factory("BsModal", ->
    (id) ->
      new ModalManager(id)
)
