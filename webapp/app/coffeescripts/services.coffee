# Services

angular.module('app.goatdServices', ['ng'])
  .factory('Rpc', ['$http'], ($http) ->
    RpcFactory = (service, actions...) ->
      Rpc = -> @service = service

      for action in actions
        do (action) ->
          Rpc.prototype[action] = (params, success = null) ->
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
      Rpc

    RpcFactory
  )
