defaultLabels =
  title: "TODO: set title"
  action: "TODO: set action"

class ModalManager
  constructor: (id, options) ->
    @sel = "##{id}"
    @save = options.save || (-> false)
    @attributes = options.attributes || []
    @labels = {}

  open: (options = {}) ->
    _.extend(@labels, defaultLabels, options.labels || {})
    @[name] = value for name, value of options.values if options.values?
    $(@sel)
      .modal("show")
      .one("shown", (event) ->
        $(event.target)
          .find("form input[type=text]:visible:first")[0]?.focus()
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
