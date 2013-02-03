# # Services
# angular.module('app.goatdServices', ['ngResource'])
#   .factory 'Overview', ($resource) ->
#     $resource('rpc', {method: "Overview.List"},
#       list: {method: 'POST', params: {id: (new Date).getTime()}, isArray: true})
