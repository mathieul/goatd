angular.module("atdServices").factory("Teammate", [
  "$resource",
  ($resource) ->
    $resource '/teammates/:uid', {uid: '@uid'},
      index:   {method: 'GET', isArray: true}
      create:  {method: 'POST'}
      update:  {method: 'PUT'}
      destroy: {method: 'DELETE'}
])
