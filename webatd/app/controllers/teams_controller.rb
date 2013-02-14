class TeamsController < ApplicationController
  include ActionController::MimeResponds

  respond_to :json

  def index
    list = nil
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("team", "list")
      list = receive_response
    end
    respond_with list
  end

  def create
    team = nil
    TcpClient.new(AppConfig[:atd_address]) do
      send_request("team", "create", name: params[:name])
      team = receive_response
    end
    respond_with team
  end
end
