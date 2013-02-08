Webatd::Application.routes.draw do
  root :to => 'home#index'

  resources :teams, only: :index

  # match ':controller(/:action(/:id))(.:format)'
end
