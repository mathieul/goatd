angular.module("atdServices").factory("Team", [
  "$resource",
  ($resource) ->
    $resource('teams')
])
