class ModalManager
  constructor: (id, options) ->
    @sel = "##{id}"
    @save = options.save || (-> false)
    @attributes = options.attributes || []

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
          attributes = {}
          for name in @attributes
            attributes[name] = this[name]
          @save(attributes)
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
