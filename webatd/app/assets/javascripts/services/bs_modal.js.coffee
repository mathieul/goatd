class ModalManager
  constructor: (id, options) ->
    @sel = "##{id}"
    @save = options.save || (-> false)

  open: (options = {}) ->
    options.title ||= "TODO: set title"
    options.action ||= "TODO: set action"
    this[name] = value for name, value of options
    $(@sel)
      .modal("show")
      .one("shown", (event) ->
        $(event.target).find("form input[type=text]:visible:first")[0].focus()
      )
      .find("form")
        .on("submit", (event) =>
          @save(name: @name)
          @close()
          event.preventDefault()
        )

  close: ->
    $(@sel)
      .modal("hide")
      .find("form")
        .off("submit")

angular.module("atdServices").factory("BsModal", ->
    (id, options = {}) ->
      new ModalManager(id, options)
)
