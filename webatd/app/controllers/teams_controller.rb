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
    team, attributes = nil, params[:team]
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("team", "create", name: attributes[:name])
      team = receive_response
    end
    render json: team, status: :created
  end

  def update
    uid, team = params[:id], params[:team].slice(:name)
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("team", "update", uid: uid, team: team)
      res = receive_response
    end
    render json: "", status: :ok
  end
end
