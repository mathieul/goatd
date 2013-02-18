Webatd::Application.routes.draw do
  root :to => 'home#index'

  match "/overview" => "home#overview"
  resources :teams, only: [:index, :create, :update, :destroy]
end
