angular.module("atdServices").factory("Team", [
  "$resource",
  ($resource) ->
    $resource '/teams/:uid', {uid: '@uid'},
      index:   {method: 'GET', isArray: true}
      create:  {method: 'POST'}
      update:  {method: 'PUT'}
      destroy: {method: 'DELETE'}
])
