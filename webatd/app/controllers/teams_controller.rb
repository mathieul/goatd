class TeamsController < ApplicationController
  def index
    list = nil
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("team", "list")
      list = receive_response
    end
    render json: list["teams"]
  end

  def create
    team, name = nil, params[:name]
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("team", "create", name: name)
      team = receive_response
    end
    render json: team, status: :created
  end
end
