# Services
angular.module("goatdServices").factory("Rpc", [
  "$http",
  ($http) ->

    serviceFactory = (service, actions...) ->
      ServiceConstructor = -> @service = service
      for action in actions
        do (action) ->
          ServiceConstructor.prototype[action] = (params, success = null) ->
            if success is null
              [params, success] = [[], params]
            else
              params = [params]

            $http
              .post("rpc",
                method: "#{@service}.#{action}"
                params: params
                id:     (new Date).getTime()
              ,
                headers:
                  "Content-Type": "application/json"
              )
              .success (data) ->
                success(data.result)
      new ServiceConstructor
    serviceFactory

  ])
